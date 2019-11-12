package main

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"hash"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

const sessionCookieName = "_happytohelp_session"

type appContext string

const (
	contextUser appContext = "contextUser"
)

type App struct {
	db      *sql.DB
	hash    hash.Hash
	baseurl string
	port    string
}

type User struct {
	Username    string
	Displayname string
	Email       string
}

func main() {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("can't open db: ", err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal("can't ping db: ", err)
	}
	app := App{
		db:      db,
		baseurl: os.Getenv("BASEURL"),
		port:    os.Getenv("PORT"),
		hash:    hmac.New(sha256.New, []byte(os.Getenv("HMAC_KEY"))),
	}
	hub := newHub()
	go hub.run()

	http.HandleFunc("/", app.root)

	// Student
	http.HandleFunc("/student/category", app.studentCategory)
	http.HandleFunc("/student/disclosure", app.studentDisclosure)
	http.HandleFunc("/student/chat", app.studentChat)
	http.HandleFunc("/student/chat/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})

	// Counsellor
	http.HandleFunc("/counsellor/login", nusRedirect(app.baseurl+app.port+"/counsellor/login/callback"))
	http.HandleFunc("/counsellor/login/callback", nusAuthenticate(app.setsession(redirect("/counsellor/category"))))
	http.HandleFunc("/counsellor/category", app.getsession(app.counsellorCategory))
	http.HandleFunc("/counsellor/choose", app.getsession(app.counsellorChoose))
	http.HandleFunc("/counsellor/chat", app.getsession(app.counsellorChat))
	http.HandleFunc("/counsellor/chat/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	fmt.Printf("Listening on " + app.baseurl + app.port)
	err = http.ListenAndServe(app.port, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func (user User) String() string {
	return fmt.Sprintf("username:%s displayname:%s email:%s", user.Username, user.Displayname, user.Email)
}

func (app *App) sign(input []byte) (output string) {
	app.hash.Reset()
	app.hash.Write(input)
	b := app.hash.Sum(nil)
	output = base64.URLEncoding.EncodeToString(b)
	return output
}

func (app *App) serialize(input interface{}) (output string, err error) {
	payload, err := json.Marshal(input)
	if err != nil {
		return "", wrap(err)
	}
	encodedPayload := base64.URLEncoding.EncodeToString(payload)
	signature := app.sign(payload)
	output = encodedPayload + "." + signature
	return output, nil
}

func (app *App) deserialize(input string, output interface{}) (err error) {
	strs := strings.SplitN(input, ".", 2)
	if len(strs) < 2 {
		return fmt.Errorf("No '.' found in input %s", input)
	}
	encodedPayload := strs[0]
	providedSignature := strs[1]
	payload, err := base64.URLEncoding.DecodeString(encodedPayload)
	if err != nil {
		return wrap(err)
	}
	computedSignature := app.sign(payload)
	if providedSignature != computedSignature {
		return fmt.Errorf("providedSignature did not match computedSignature %+v", struct {
			Provided string
			Computed string
		}{providedSignature, computedSignature})
	}
	err = json.Unmarshal(payload, output)
	return wrap(err)
}

func notfound(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "404.html")
}

func (app *App) root(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		notfound(w, r)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	err := render(w, r, nil, "landing.html")
	if err != nil {
		dump(w, err)
	}
}
func render(w http.ResponseWriter, r *http.Request, data interface{}, filenames ...string) (err error) {
	if len(filenames) == 0 {
		return wrap(fmt.Errorf("no filenames provided to Render"))
	}
	funcs := template.FuncMap{}
	funcs["AppGetUser"] = func(r *http.Request) func() User {
		user, _ := r.Context().Value(contextUser).(User)
		return func() User {
			return user
		}
	}(r)
	filenames = append(filenames, "navbar.html")
	t, err := template.New(filepath.Base(filenames[0])).Funcs(funcs).ParseFiles(filenames...)
	if err != nil {
		return wrap(err)
	}
	w.Header().Set("Content-Type", "text/html")
	err = t.Execute(w, data)
	return wrap(err)
}

func wrap(err error) error {
	if err != nil {
		if err == sql.ErrNoRows || err == http.ErrNoCookie {
			return err
		} else {
			pc, filename, linenr, _ := runtime.Caller(1)
			return errors.Wrapf(err, " • error in function[%s] file[%s] line[%d]", runtime.FuncForPC(pc).Name(), filename, linenr)
		}
	} else {
		return nil
	}
}

func dump(w io.Writer, err error) {
	pc, filename, linenr, _ := runtime.Caller(1)
	err = errors.Wrapf(err, " • error in function[%s] file[%s] line[%d]", runtime.FuncForPC(pc).Name(), filename, linenr)
	fmtedErr := strings.Replace(err.Error(), " • ", "\n\n", -1)
	fmt.Fprintf(w, fmtedErr)
}

func (app *App) setsession(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, ok1 := r.Context().Value("username").(string)
		displayname, ok2 := r.Context().Value("displayname").(string)
		email, ok3 := r.Context().Value("email").(string)
		user := User{Username: username, Displayname: displayname, Email: email}
		if !ok1 || !ok2 || !ok3 {
			fmt.Fprintf(w, "could not get all details %s", user)
			return
		}
		serializedUser, err := app.serialize(user)
		if err != nil {
			dump(w, err)
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:     sessionCookieName,
			Value:    serializedUser,
			HttpOnly: true, // disable JavaScript access to cookie
			// Secure:   true, // allow sending only over HTTPS
			SameSite: http.SameSiteLaxMode,
			MaxAge:   int((time.Hour * 24 * 30).Seconds()), // one month
			Path:     "/",
		})
		next(w, r)
	}
}

func (app *App) getsession(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionCookie, err := r.Cookie(sessionCookieName)
		if err != nil {
			http.Error(w, "", http.StatusUnauthorized)
			return
		}
		var user User
		err = app.deserialize(sessionCookie.Value, &user)
		if err != nil {
			dump(w, err)
			return
		}
		ctx := r.Context()
		ctx = context.WithValue(ctx, contextUser, user)
		next(w, r.WithContext(ctx))
	}
}

func redirect(url string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, url, http.StatusMovedPermanently)
	}
}
