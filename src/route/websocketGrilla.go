package route

import (
	"logger"
	"mongo"

	"fmt"
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}
var thread *mongo.Thread

// connection is an middleman between the websocket connection and the hub.
type connection struct {
	// The websocket connection.
	ws *websocket.Conn

	// Buffered channel of outbound messages.
	send chan *mongo.Thread
}

// readPump pumps messages from the websocket connection to the hub.
func (c *connection) readPump() {
	defer func() {
		H.Unregister <- c
		c.ws.Close()
	}()
	c.ws.SetReadLimit(maxMessageSize)
	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetPongHandler(func(string) error { c.ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {

		err := c.ws.ReadJSON(&thread)
		if err != nil {
			if thread == nil {
				break
			}
			logger.Output(
				logrus.Fields{
					logger.ErrMsg:   fmt.Errorf("%v", err),
					logger.TraceMsg: logger.Trace(),
					"Thread":        fmt.Sprintf("%+v\n", thread),
				},
				"can not read JSON or Closed Websocket",
				logger.Debug,
			)
			break
		}

		if err := mongo.UpsertNewTimeLine(thread); err != nil {
			logger.Output(
				logrus.Fields{
					logger.ErrMsg:   fmt.Errorf("%v", err),
					logger.TraceMsg: logger.Trace(),
					"Thread":        fmt.Sprintf("%+v\n", thread),
				},
				"can not upsert thread",
				logger.Error,
			)
		}

		H.Broadcast <- thread
	}
}

// write writes a message with the given message type and payload.
func (c *connection) write(mt int, payload []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteMessage(mt, payload)
}

// writePump pumps messages from the hub to the websocket connection.
func (c *connection) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.ws.Close()
	}()
	for {
		select {
		case threadToResponse, ok := <-c.send:
			if !ok {
				c.write(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.ws.WriteJSON(threadToResponse); err != nil {
				logger.Output(
					logrus.Fields{
						logger.ErrMsg:        err,
						logger.TraceMsg:      logger.Trace(),
						"Thread To Response": threadToResponse,
					},
					"can not write thread",
					logger.Error,
				)
				return
			}
		case <-ticker.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				logger.Output(
					logrus.Fields{
						logger.ErrMsg:   err,
						logger.TraceMsg: logger.Trace(),
						"Ping Message":  websocket.PingMessage,
					},
					"can not ping message",
					logger.Error,
				)
				return
			}
		}
	}
}

// serveWs handles websocket requests from the peer.
func ServeWs(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Output(
			logrus.Fields{
				logger.ErrMsg:         err,
				logger.TraceMsg:       logger.Trace(),
				"http.ResponseWriter": w,
				"http.Request":        r,
			},
			"can not make websocket",
			logger.Error,
		)
		return
	}
	c := &connection{send: make(chan *mongo.Thread, 256), ws: ws}
	H.Register <- c
	go c.writePump()
	c.readPump()
}
