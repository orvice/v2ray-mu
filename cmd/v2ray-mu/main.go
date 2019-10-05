package main

import (
	"fmt"
	"os"

	"github.com/orvice/v2ray-mu/server"
	"github.com/weeon/utils/process"
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

	process.WaitSignal()
}
