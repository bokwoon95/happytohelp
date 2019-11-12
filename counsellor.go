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
		app.notfound(w, r)
	}
}

func (app *App) counsellorTopic(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	user, ok := r.Context().Value(contextUser).(User)
	if !ok {
		fmt.Fprintf(w, "Unable to get user from context %s", user)
		return
	}
	switch r.Method {
	case "GET":
		err := render(w, r, nil, "counsellor_topics.html")
		if err != nil {
			dump(w, err)
		}
	default:
		app.notfound(w, r)
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
		app.chatrooms.Lock()
		defer app.chatrooms.Unlock()
		type TupleRoom struct {
			Key  string
			Room *Chatroom
		}
		var pendingRooms []TupleRoom
		for key, room := range app.chatrooms.pendingRooms {
			pendingRooms = append(pendingRooms, TupleRoom{Key: key, Room: room})
		}
		type Data struct {
			PendingRooms []TupleRoom
		}
		data := Data{PendingRooms: pendingRooms}
		err := render(w, r, data, "counsellor_choose.html")
		if err != nil {
			dump(w, err)
		}
	default:
		app.notfound(w, r)
	}
}
