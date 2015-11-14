package route

import (
	l "logger"
	m "mongo"

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
var thread *m.Thread
var newThread *m.Thread

// connection is an middleman between the websocket connection and the hub.
type connection struct {
	// The websocket connection.
	ws *websocket.Conn

	// Buffered channel of outbound messages.
	send chan *m.Thread
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
			l.Output(
				logrus.Fields{
					l.ErrMsg:   l.Errorf(err),
					l.TraceMsg: l.Trace(),
					"Thread":   l.Sprintf(thread),
				},
				"can not read JSON or Closed Websocket",
				l.Debug,
			)
			break
		}

		if newThread, err = m.UpdateExistingTimeLine(thread); err != nil {
			l.PutErr(err, l.Trace(), l.E_R_Upsert, thread)
		}

		H.Broadcast <- newThread
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
				l.Output(
					logrus.Fields{
						l.ErrMsg:   l.Errorf(err),
						l.TraceMsg: l.Trace(),
						"Thread":   l.Sprintf(threadToResponse),
					},
					"can not write JSON",
					l.Debug,
				)
				return
			}
		case <-ticker.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				l.PutErr(err, l.Trace(), l.E_R_PingMsg, websocket.PingMessage)
				return
			}
		}
	}
}

// serveWs handles websocket requests from the peer.
func ServeWs(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		l.PutErr(err, l.Trace(), l.E_R_Upgrader, nil)
		return
	}
	c := &connection{send: make(chan *m.Thread, 256), ws: ws}
	H.Register <- c
	go c.writePump()
	c.readPump()
}
