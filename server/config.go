package server

import (
	"time"

	"github.com/weeon/utils"
)

var (
	cfg = new(Config)
)

type Config struct {
	WebApi   WebApiCfg
	Base     BaseCfg
	SyncTime time.Duration

	TrojanApiServerAddr string

	V2rayClientAddr string
	V2rayTag        string

	LogPath string
}

type BaseCfg struct {
}

type WebApiCfg struct {
	Url    string
	Token  string
	NodeID string
}

func initCfg() {
	cfg.WebApi = WebApiCfg{
		Url:    utils.GetEnv("MU_URI"),
		Token:  utils.GetEnv("MU_TOKEN"),
		NodeID: utils.GetEnv("MU_NODE_ID"),
	}
	st := utils.GetEnvInt("SYNC_TIME", 60)
	cfg.SyncTime = time.Second * time.Duration(st)
	cfg.TrojanApiServerAddr = utils.GetEnv("TROJAN_API_SERVER_ADDR")
	cfg.V2rayClientAddr = utils.GetEnv("V2RAY_ADDR")
	cfg.V2rayTag = utils.GetEnv("V2RAY_TAG")
	cfg.LogPath = utils.GetEnv("LOG_DIR", "/var/log/v2ray-mu/")
}
