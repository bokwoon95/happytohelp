package main

import (
	"log"
	"net/http"
)

func (app App) counsellorChat(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	switch r.Method {
	case "GET":
		err := render(w, r, nil, "counsellor_chat.html")
		if err != nil {
			dump(w, err)
		}
	default:
		notfound(w, r)
	}
}

func (app App) counsellorCategory(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	switch r.Method {
	case "GET":
		err := render(w, r, nil, "student_category.html")
		if err != nil {
			dump(w, err)
		}
	default:
		notfound(w, r)
	}
}
