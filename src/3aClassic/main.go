package main

import (
	"fmt"
	//	"golang.org/x/net/websocket"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"route"
)

func main() {
	route.Register()
	go route.H.Run()
	//	http.Handle("/ws/timeLine", websocket.Handler(route.EchoHandler))
	http.HandleFunc("/ws/timeLine", route.ServeWs)
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
