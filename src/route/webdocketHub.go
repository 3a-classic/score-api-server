// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package route

import (
	"log"
)

// hub maintains the set of active connections and broadcasts messages to the
// connections.
type Hub struct {
	// Registered connections.
	Connections map[*connection]bool

	// Inbound messages from the connections.
	Broadcast chan []byte

	// Register requests from the connections.
	Register chan *connection

	// Unregister requests from connections.
	Unregister chan *connection
}

var H = Hub{
	Broadcast:   make(chan []byte),
	Register:    make(chan *connection),
	Unregister:  make(chan *connection),
	Connections: make(map[*connection]bool),
}

func (h *Hub) Run() {
	for {
		select {
		case c := <-h.Register:
			h.Connections[c] = true
			log.Println("register", h.Connections[c])
		case c := <-h.Unregister:
			if _, ok := h.Connections[c]; ok {
				log.Println("unregister", h.Connections[c])
				delete(h.Connections, c)
				close(c.send)
			}
		case m := <-h.Broadcast:
			for c := range h.Connections {
				log.Println("broadcaset", c)
				select {
				case c.send <- m:
				default:
					close(c.send)
					delete(h.Connections, c)
				}
			}
		}
	}
}
