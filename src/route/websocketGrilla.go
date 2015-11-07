package route

import (
	//	"golang.org/x/net/websocket"
	//	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"mongo"
	"net/http"
	//	"runtime"
	"time"
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
		//		_, message, err := c.ws.ReadMessage()
		err := c.ws.ReadJSON(&thread)
		//		fmt.Printf("%+v\n", thread)
		if err != nil {
			log.Println("JSON Read ERR", err)
			//	break
		}
		//		H.Broadcast <- message

		if err := mongo.UpsertNewTimeLine(thread); err != nil {
			log.Println("reaction update err", err)
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
		case threadToWrite, ok := <-c.send:
			if !ok {
				c.write(websocket.CloseMessage, []byte{})
				return
			}
			//			if err := c.write(websocket.TextMessage, message); err != nil {
			if err := c.ws.WriteJSON(threadToWrite); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

// serveWs handles websocket requests from the peer.
func ServeWs(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	//	c := &connection{send: make(chan []byte, 256), ws: ws}
	c := &connection{send: make(chan *mongo.Thread, 256), ws: ws}
	H.Register <- c
	go c.writePump()
	c.readPump()
}

//var (
//	ActiveClients = make(map[ClientConn]int)
//)
//
//type ClientConn struct {
//	websocket *websocket.Conn
//	clientIP  string
//}
//
//func EchoHandler(ws *websocket.Conn) {
//
//	defer func() {
//		if err := ws.Close(); err != nil {
//			log.Println("Websocket could not be closed", err.Error())
//		}
//	}()
//	client := ws.Request().RemoteAddr
//	log.Println("Client connected:", client)
//	sockCli := ClientConn{ws, client}
//
//	ActiveClients[sockCli] = 0
//	log.Println("Number of clients connected ...", len(ActiveClients))
//
//	var thread mongo.Thread
//	for {
//		log.Println("goroutine num", runtime.NumGoroutine())
//
//		if err := websocket.JSON.Receive(ws, &thread); err != nil {
//			log.Println("Websocket Disconnected waiting", err.Error())
//			delete(ActiveClients, sockCli)
//			log.Println("Number of clients still connected ...", len(ActiveClients))
//			return
//		}
//
//		if err := mongo.UpsertNewTimeLine(&thread); err != nil {
//			log.Println("cannot insert data to mongo", err.Error())
//		}
//
//		log.Println("Number of channel ", len(mongo.FinChan))
//		for cs, _ := range ActiveClients {
//			//		if err = Message.Send(cs.websocket, clientMessage); err != nil {
//			if err := websocket.JSON.Send(cs.websocket, thread); err != nil {
//				// we could not send the message to a peer
//				log.Println("Could not send message to ", cs.clientIP, err.Error())
//			}
//		}
//	}
//}
