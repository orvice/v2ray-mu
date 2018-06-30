package main

import (
	"fmt"
	"github.com/orvice/kit/log"
	"os"
	"os/signal"
	"syscall"
)

var (
	logger log.Logger
)

func init() {
}

func main() {
	var err error
	initCfg()
	logger = log.NewFileLogger(cfg.LogPath)

	InitWebApi()

	um, err := NewUserManager()
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	go um.Run()

	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, os.Interrupt, os.Kill, syscall.SIGTERM)

	<-osSignals
}
