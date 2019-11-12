package main

import (
	"log"
	"net/http"
)

func (app App) studentCategory(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	err := render(w, r, nil, "student_category.html")
	if err != nil {
		dump(w, err)
	}
}

func (app App) studentDisclosure(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	switch r.Method {
	case "GET":
		err := render(w, r, nil, "student_disclosure.html")
		if err != nil {
			dump(w, err)
		}
	case "POST":
		http.Redirect(w, r, "/student/chat", http.StatusSeeOther)
	default:
		notfound(w, r)
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
