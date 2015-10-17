package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"route"
)

func main() {
	route.Register()
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
