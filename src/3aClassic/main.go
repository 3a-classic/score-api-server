package main

import (
	"logger"
	"route"

	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/Sirupsen/logrus"
)

func main() {
	route.Register()
	go route.H.Run()
	http.HandleFunc("/ws/timeLine", route.ServeWs)
	http.ListenAndServe(":80", nil)
	shutdownHook()
}

func shutdownHook() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
	logger.Output(
		logrus.Fields{},
		"Shutdown",
		logger.Info,
	)
	os.Exit(0)
}
