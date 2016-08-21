package main

import (
	l "github.com/3a-classic/score-api-server/logger"
	r "github.com/3a-classic/score-api-server/route"

	"os"
	"os/signal"
	"syscall"
)

func main() {
	l.Output(nil, "Start API server", l.Info)
	r.Register()
	shutdownHook()
}

func shutdownHook() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
	l.Output(nil, "Shutdown", l.Info)
	os.Exit(0)
}
