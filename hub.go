package main

import "sync"

// Chatroom maintains the set of active clients and broadcasts messages to the
// clients.
type Chatroom struct {
	sync.Mutex

	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	// Shutdown chatroom
	shutdown chan struct{}
}

func newChatroom() *Chatroom {
	return &Chatroom{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		shutdown:   make(chan struct{}),
	}
}

func (room *Chatroom) run() {
	for {
		select {
		case <-room.shutdown:
			return
		case client := <-room.register:
			room.Lock()
			room.clients[client] = true
			room.Unlock()
		case client := <-room.unregister:
			room.Lock()
			if _, ok := room.clients[client]; ok {
				delete(room.clients, client)
				close(client.send)
			}
			room.Unlock()
		case message := <-room.broadcast:
			room.Lock()
			for client := range room.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(room.clients, client)
				}
			}
			room.Unlock()
		}
	}
}
