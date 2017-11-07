package main

import (
	"github.com/orvice/utils/env"
	"time"
)

var (
	cfg = new(Config)
)

type Config struct {
	WebApi   WebApiCfg
	Base     BaseCfg
	SyncTime time.Duration
}

type BaseCfg struct {
}

type WebApiCfg struct {
	Url    string
	Token  string
	NodeId int
}

func initCfg() {
	cfg.WebApi = WebApiCfg{
		Url:    env.Get("MU_URI"),
		Token:  env.Get("MU_TOKEN"),
		NodeId: env.GetInt("MU_NODE_ID"),
	}
	st := env.GetInt("SYNC_TIME", 60)
	cfg.SyncTime = time.Second * time.Duration(st)
}
