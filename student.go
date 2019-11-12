package main

import (
	"log"
	"net/http"
	"net/url"
)

func (app *App) studentTopic(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	err := render(w, r, nil, "student_topics.html")
	if err != nil {
		dump(w, err)
	}
}

func (app *App) studentDisclosure(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	switch r.Method {
	case "GET":
		err := render(w, r, nil, "student_disclosure.html")
		if err != nil {
			dump(w, err)
		}
	case "POST":
		r.ParseForm()
		topics := r.Form["topics"]
		u, err := url.Parse("/student/chat")
		if err != nil {
			dump(w, err)
			return
		}
		q := u.Query()
		for _, topic := range topics {
			q.Add("topics", topic)
		}
		u.RawQuery = q.Encode()
		http.Redirect(w, r, u.String(), http.StatusSeeOther)
	default:
		app.notfound(w, r)
	}
}

func (app *App) studentChat(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	switch r.Method {
	case "GET":
		err := render(w, r, nil, "student_chat.html")
		if err != nil {
			dump(w, err)
		}
	default:
		app.notfound(w, r)
		return
	}
}
