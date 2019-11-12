package main

import (
	"fmt"
	"log"
	"net/http"
)

func (app *App) counsellorChat(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	user, ok := r.Context().Value(contextUser).(User)
	if !ok {
		fmt.Fprintf(w, "Unable to get user from context %s", user)
		return
	}
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

func (app *App) counsellorCategory(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	user, ok := r.Context().Value(contextUser).(User)
	if !ok {
		fmt.Fprintf(w, "Unable to get user from context %s", user)
		return
	}
	switch r.Method {
	case "GET":
		err := render(w, r, nil, "counsellor_category.html")
		if err != nil {
			dump(w, err)
		}
	default:
		notfound(w, r)
	}
}

func (app *App) counsellorChoose(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	user, ok := r.Context().Value(contextUser).(User)
	if !ok {
		fmt.Fprintf(w, "Unable to get user from context %s", user)
		return
	}
	switch r.Method {
	case "GET":
		err := render(w, r, nil, "counsellor_choose.html")
		if err != nil {
			dump(w, err)
		}
	default:
		notfound(w, r)
	}
}
