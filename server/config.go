package server

import (
	"time"

	"github.com/orvice/utils/env"
)

var (
	cfg = new(Config)
)

type Config struct {
	WebApi   WebApiCfg
	Base     BaseCfg
	SyncTime time.Duration

	V2rayClientAddr string
	V2rayTag        string

	LogPath string
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
	cfg.V2rayClientAddr = env.Get("V2RAY_ADDR")
	cfg.V2rayTag = env.Get("V2RAY_TAG")
	cfg.LogPath = env.Get("LOG_PATH", "/var/log/v2ray-mu.log")
}
