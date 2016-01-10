package main

import (
	l "logger"
	r "route"

	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	r.Register()
	go r.H.Run()
	http.HandleFunc("/ws/timeLine", r.ServeWs)
	http.ListenAndServe(":80", nil)
	shutdownHook()
}

func shutdownHook() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
	l.Output(nil, "Shutdown", l.Info)
	os.Exit(0)
}
