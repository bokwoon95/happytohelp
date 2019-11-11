package main

import (
	"html/template"
	"log"
	"net/http"
)

func notfound(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "404.html")
}

func landing(w http.ResponseWriter, r *http.Request) {
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

func category(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "category.html")
}

func disclosure(w http.ResponseWriter, r *http.Request) {
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

func studentChat(w http.ResponseWriter, r *http.Request) {
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

func counsellorChat(w http.ResponseWriter, r *http.Request) {
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
	hub := newHub()
	go hub.run()
	http.HandleFunc("/", landing)
	http.HandleFunc("/category", category)
	http.HandleFunc("/disclosure", disclosure)
	http.HandleFunc("/student-chat", studentChat)
	http.HandleFunc("/student-chat/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})
	http.HandleFunc("/counsellor-chat", counsellorChat)
	http.HandleFunc("/counsellor-chat/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
