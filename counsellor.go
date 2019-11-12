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
	key := r.FormValue("key")
	app.chatrooms.Lock()
	defer app.chatrooms.Unlock()
	_, ok = app.chatrooms.pendingRooms[key]
	if !ok {
		http.Error(w, fmt.Sprintf("Invalid room key %s", key), http.StatusBadRequest)
		return
	}
	switch r.Method {
	case "GET", "POST":
		log.Println("gotcha")
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
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
		type KeyRoom struct {
			Key  string
			Room *Chatroom
		}
		var pendingRooms []KeyRoom
		for key, room := range app.chatrooms.pendingRooms {
			pendingRooms = append(pendingRooms, KeyRoom{Key: key, Room: room})
		}
		type Data struct {
			PendingRooms []KeyRoom
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

func (app *App) counsellorWs(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	user, ok := r.Context().Value(contextUser).(User)
	if !ok {
		fmt.Fprintf(w, "Unable to get user from context %s", user)
		return
	}
	key := r.FormValue("key")
	app.chatrooms.Lock()
	defer app.chatrooms.Unlock()
	room, ok := app.chatrooms.pendingRooms[key]
	room.broadcast <- []byte(fmt.Sprintf("Counsellor found. Say hi to %s!", user.Displayname))
	if !ok {
		http.Error(w, fmt.Sprintf("Invalid room key %s", key), http.StatusBadRequest)
		return
	}
	serveWs(room, w, r)
}
