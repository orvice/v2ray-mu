package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"v2ray.com/core"

	"bytes"
	"time"
	_ "v2ray.com/core/main/distro/all"
	_ "v2ray.com/core/tools/conf"
)

var (
	server core.Server

	ApiAddr string = ""
)

func init() {
}

func GetConfigFormat() core.ConfigFormat {
	return core.ConfigFormat_JSON
}

func startV2Ray() (core.Server, error) {
	res, err := apiClient.GetV2rayUsersData()
	if err != nil {
		return nil, err
	}
	configInput := bytes.NewBufferString(res)
	config, err := core.LoadConfig(GetConfigFormat(), configInput)
	if err != nil {
		return nil, newError("failed to read config from json ").Base(err)
	}

	server, err := core.New(config)
	if err != nil {
		return nil, newError("failed to create initialize").Base(err)
	}

	return server, nil
}

func run() {
	var err error
	server, err = startV2Ray()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if err := server.Start(); err != nil {
		fmt.Println("Failed to start", err)
	}
}

func check() {
	var lastRes string
	var err error
	lastRes, err = apiClient.GetV2rayUsersData()
	if err != nil {
		panic(err)
	}
	time.Sleep(cfg.SyncTime)
	for {
		res, err := apiClient.GetV2rayUsersData()
		if err != nil {
			fmt.Println(err)
			time.Sleep(cfg.SyncTime)
			continue
		}
		if lastRes != res {
			fmt.Println("restart server")
			server.Close()
			go run()
		}
		time.Sleep(cfg.SyncTime)
	}
}

func main() {
	go pprof()
	core.PrintVersion()
	initCfg()
	InitWebApi()

	go run()
	go check()

	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, os.Interrupt, os.Kill, syscall.SIGTERM)

	<-osSignals
	server.Close()
}
