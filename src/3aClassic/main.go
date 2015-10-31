package main

import (
	"fmt"
	"golang.org/x/net/websocket"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"route"
)

func main() {
	route.Register()
	http.Handle("/ws/timeLine", websocket.Handler(route.EchoHandler))
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
