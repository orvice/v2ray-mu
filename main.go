package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

var (
	ApiAddr string = ""
)

func init() {
}

func main() {
	var err error
	initCfg()
	InitWebApi()
	InitUserManager()

	err = InitV2rayManager()
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	go daemon()

	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, os.Interrupt, os.Kill, syscall.SIGTERM)

	<-osSignals
}
