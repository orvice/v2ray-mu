package main

import (
	"fmt"
	"log"
	"os"

	"github.com/google/gops/agent"
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

	for _, v := range um {
		go func(v *server.UserManager) {
			_ = v.Run()
		}(v)
	}

	go func() {
		if err := agent.Listen(agent.Options{}); err != nil {
			log.Fatal(err)
		}
	}()

	process.WaitSignal()
}
