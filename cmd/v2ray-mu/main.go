package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/orvice/v2ray-mu/server"
)

func main() {
	server.Init()

	server.InitWebApi()

	um, err := server.NewUserManager()
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	go um.Run()

	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, os.Interrupt, os.Kill, syscall.SIGTERM)

	<-osSignals
}
