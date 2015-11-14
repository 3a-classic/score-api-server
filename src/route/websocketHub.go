// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package route

import (
	"logger"
	"mongo"

	"github.com/Sirupsen/logrus"
)

// hub maintains the set of active connections and broadcasts messages to the
// connections.
type Hub struct {
	// Registered connections.
	Connections map[*connection]bool

	// Inbound messages from the connections.
	Broadcast chan *mongo.Thread

	// Register requests from the connections.
	Register chan *connection

	// Unregister requests from connections.
	Unregister chan *connection
}

var H = Hub{
	Broadcast:   make(chan *mongo.Thread),
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
			logger.Output(
				logrus.Fields{
					"New Connection": c,
					"Connections":    h.Connections,
				},
				"Register websocket",
				logger.Debug,
			)
		case c := <-h.Unregister:
			if _, ok := h.Connections[c]; ok {
				logger.Output(
					logrus.Fields{
						"New Connection": c,
						"Connections":    h.Connections,
					},
					"Unregister websocket",
					logger.Debug,
				)
				delete(h.Connections, c)
				close(c.send)
			}
		case m := <-h.Broadcast:
			logger.Output(
				logrus.Fields{
					"Broad Cast":  m,
					"Connections": h.Connections,
				},
				"Broad Cast websocket",
				logger.Debug,
			)
			for c := range h.Connections {
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

func newThreadCheck() {
	for {
		select {
		case t := <-mongo.ThreadChan:
			logger.Output(
				logrus.Fields{
					"Threa Chain": t,
				},
				"New Thread websocket",
				logger.Debug,
			)

			H.Broadcast <- t
		}
	}
}
