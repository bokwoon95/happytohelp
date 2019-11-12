package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
)

func (app *App) studentTopics(w http.ResponseWriter, r *http.Request) {
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
		// Generate chat room url
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
		key, err := generateRandomString()
		if err != nil {
			dump(w, err)
			return
		}
		q.Add("room", key)
		u.RawQuery = q.Encode()
		// Create new chat room
		app.chatrooms.Lock()
		defer app.chatrooms.Unlock()
		room := newChatroom()
		go room.run()
		room.Topics = topics
		room.PinnedMessage = r.FormValue("disclosure")
		app.chatrooms.pendingRooms[key] = room
		http.Redirect(w, r, u.String(), http.StatusSeeOther)
	default:
		app.notfound(w, r)
	}
}

func (app *App) studentChat(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	switch r.Method {
	case "GET":
		key := r.FormValue("room")
		app.chatrooms.Lock()
		defer app.chatrooms.Unlock()
		_, ok := app.chatrooms.pendingRooms[key]
		if !ok {
			http.Error(w, fmt.Sprintf("Invalid room key %s", key), http.StatusBadRequest)
			return
		}
		err := render(w, r, nil, "student_chat.html")
		if err != nil {
			dump(w, err)
		}
	default:
		app.notfound(w, r)
		return
	}
}

func (app *App) studentWs(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	key := r.FormValue("room")
	app.chatrooms.Lock()
	defer app.chatrooms.Unlock()
	room, ok := app.chatrooms.pendingRooms[key]
	if !ok {
		http.Error(w, fmt.Sprintf("Invalid room key %s", key), http.StatusBadRequest)
		return
	}
	serveWs(room, w, r)
}
