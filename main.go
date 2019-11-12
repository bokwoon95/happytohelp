package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

type App struct {
	db *sql.DB
}

func notfound(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "404.html")
}

func (app App) root(w http.ResponseWriter, r *http.Request) {
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

func (app App) category(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	err := render(w, r, nil, "category.html")
	if err != nil {
		dump(w, err)
	}
}

func (app App) disclosure(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	switch r.Method {
	case "GET":
		err := render(w, r, nil, "disclosure.html")
		if err != nil {
			dump(w, err)
			return
		}
	case "POST":
		http.Redirect(w, r, "/student-chat", http.StatusSeeOther)
	}
}

func (app App) studentChat(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	switch r.Method {
	case "GET":
		err := render(w, r, nil, "student_chat.html")
		if err != nil {
			dump(w, err)
		}
	default:
		notfound(w, r)
		return
	}
}

func (app App) counsellorChat(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	switch r.Method {
	case "GET":
		err := render(w, r, nil, "counsellor_chat.html")
		if err != nil {
			dump(w, err)
			return
		}
	default:
		notfound(w, r)
		return
	}
}

func render(w http.ResponseWriter, r *http.Request, data interface{}, filenames ...string) (err error) {
	if len(filenames) == 0 {
		return wrap(fmt.Errorf("no filenames provided to Render"))
	}
	funcs := template.FuncMap{}
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
			// If its either a no sql row error or no cookie error, don't wrap the error otherwise it wouldn't be identifieable as sql.ErrNoRows or http.ErrNoCookie
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

func main() {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("can't open db: ", err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal("can't ping db: ", err)
	}
	app := App{db}
	hub := newHub()
	go hub.run()
	http.HandleFunc("/", app.root)
	http.HandleFunc("/category", app.category)
	http.HandleFunc("/disclosure", app.disclosure)
	http.HandleFunc("/student-chat", app.studentChat)
	http.HandleFunc("/student-chat/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})
	http.HandleFunc("/counsellor-chat", app.counsellorChat)
	http.HandleFunc("/counsellor-chat/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
