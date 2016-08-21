// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package route

import (
	l "logger"
	m "mongo"

	"github.com/Sirupsen/logrus"
)

// hub maintains the set of active connections and broadcasts messages to the
// connections.
type Hub struct {
	// Registered connections.
	Connections map[*connection]bool

	// Inbound messages from the connections.
	Broadcast chan *m.Thread

	// Register requests from the connections.
	Register chan *connection

	// Unregister requests from connections.
	Unregister chan *connection
}

var H = Hub{
	Broadcast:   make(chan *m.Thread),
	Register:    make(chan *connection),
	Unregister:  make(chan *connection),
	Connections: make(map[*connection]bool),
}

func (h *Hub) Run() {
	go newThreadCheck()

	for {
		select {
		case c := <-h.Register:
			h.Connections[c] = true
			l.Output(
				logrus.Fields{
					"New Connection": l.Sprintf(c),
					"Connections":    l.Sprintf(h.Connections),
				},
				"Register websocket",
				l.Debug,
			)
		case c := <-h.Unregister:
			if _, ok := h.Connections[c]; ok {
				l.Output(
					logrus.Fields{
						"New Connection": l.Sprintf(c),
						"Connections":    l.Sprintf(h.Connections),
					},
					"Unregister websocket",
					l.Debug,
				)
				delete(h.Connections, c)
				close(c.send)
			}
		case me := <-h.Broadcast:
			l.Output(
				logrus.Fields{
					"Broad Cast":  l.Sprintf(me),
					"Connections": l.Sprintf(h.Connections),
				},
				"Broad Cast websocket",
				l.Debug,
			)
			for c := range h.Connections {
				select {
				case c.send <- me:
				default:
					close(c.send)
					delete(h.Connections, c)
				}
			}
		}
	}
}

func newThreadCheck() {
	for {
		select {
		case t := <-m.ThreadChan:
			l.Output(
				logrus.Fields{"Threa Chain": l.Sprintf(t)},
				"New Thread websocket",
				l.Debug,
			)

			H.Broadcast <- t
		}
	}
}
