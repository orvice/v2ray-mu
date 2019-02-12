package server

import (
	"github.com/catpie/musdk-go"
	"github.com/orvice/kit/log"
)

var (
	apiClient *musdk.Client
)

var (
	logger log.Logger
)

func Init() {
	initCfg()
	logger = log.NewDefaultLogger()
}

func InitWebApi() {
	logger.Info("init mu api")
	cfg := cfg.WebApi
	apiClient = musdk.NewClient(cfg.Url, cfg.Token, cfg.NodeId, musdk.TypeV2ray)
	apiClient.SetLogger(logger)
	go apiClient.UpdateTrafficDaemon()
	return
}
