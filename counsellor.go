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
	key := r.FormValue("room")
	app.chatrooms.Lock()
	defer app.chatrooms.Unlock()
	_, ok = app.chatrooms.pendingRooms[key]
	if !ok {
		http.Error(w, fmt.Sprintf("Invalid room key %s", key), http.StatusBadRequest)
		return
	}
	switch r.Method {
	case "GET", "POST":
		type Data struct {
			User User
		}
		data := Data{User: user}
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		err := render(w, r, data, "counsellor_chat.html")
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
	// event := Event{
	// 	Sender:  user.Displayname,
	// 	Message: fmt.Sprintf("Counsellor found. Say hi to %s!", user.Displayname),
	// }
	// serializedEvent, err := json.Marshal(event)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }
	key := r.FormValue("room")
	app.chatrooms.Lock()
	defer app.chatrooms.Unlock()
	room, ok := app.chatrooms.pendingRooms[key]
	if !ok {
		http.Error(w, fmt.Sprintf("Invalid room key %s", key), http.StatusBadRequest)
		return
	}
	// room.broadcast <- []byte(serializedEvent)
	serveWs(room, w, r)
}
