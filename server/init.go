package server

import (
	"github.com/catpie/musdk-go"
	"github.com/orvice/kit/log"
)

var (
	apiClient *musdk.Client
)

var (
	logger    log.Logger
	tl        log.Logger // traffic logger
	sdkLogger log.Logger
)

func Init() {
	initCfg()
	logger = log.NewFileLogger(cfg.LogPath + "mu.log")
	tl = log.NewFileLogger(cfg.LogPath + "traffic.log")
	sdkLogger = log.NewFileLogger(cfg.LogPath + "sdk.log")
}

func InitWebApi() {
	logger.Info("init mu api")
	cfg := cfg.WebApi
	apiClient = musdk.NewClient(cfg.Url, cfg.Token, cfg.NodeId, musdk.TypeV2ray)
	apiClient.SetLogger(sdkLogger)
	go apiClient.UpdateTrafficDaemon()
	return
}
