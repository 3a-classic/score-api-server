package main

import (
	"fmt"
	"golang.org/x/net/websocket"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"route"
)

func echoHandler(ws *websocket.Conn) {

	type T struct {
		Msg string
	}

	// receive JSON type T
	//	io.Copy(ws, ws)
	var data T
	//	log.Printf("ws", ws)
	//	fmt.Println(ws)
	websocket.JSON.Receive(ws, &data)

	var msg []byte
	ws.Read(msg)
	fmt.Println(msg)
	log.Printf("data=%#v\n", data)

	// send JSON type T
	websocket.JSON.Send(ws, data)
}

func main() {
	route.Register()
	http.Handle("/ws", websocket.Handler(echoHandler))
	http.ListenAndServe(":80", nil)
	shutdownHook()
}

func shutdownHook() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
	fmt.Println("bye")
	os.Exit(0)
}
