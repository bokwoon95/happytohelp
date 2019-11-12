package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"os"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
)

type App struct {
	db *sql.DB
}

func notfound(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "404.html")
}

func (app App) landing(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		notfound(w, r)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "landing.html")
}

func (app App) category(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "category.html")
}

func (app App) disclosure(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	switch r.Method {
	case "GET":
		t, err := template.ParseFiles("disclosure.html")
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		err = t.Execute(w, nil)
		if err != nil {
			w.Write([]byte(err.Error()))
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
		t, err := template.ParseFiles("student_chat.html")
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		err = t.Execute(w, nil)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
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
		t, err := template.ParseFiles("counsellor_chat.html")
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		err = t.Execute(w, nil)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
	default:
		notfound(w, r)
		return
	}
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
	http.HandleFunc("/", app.landing)
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
