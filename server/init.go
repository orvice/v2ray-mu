package server

import (
	"github.com/catpie/musdk-go"
	"github.com/weeon/contract"
	"github.com/weeon/log"
	"go.uber.org/zap/zapcore"
)

var (
	apiClient *musdk.Client
)

var (
	logger    contract.Logger
	tl        contract.Logger // traffic logger
	sdkLogger contract.Logger
	tjLogger  contract.Logger
)

func Init() {
	initCfg()
	logger, _ = log.NewLogger(cfg.LogPath+"mu.log", zapcore.DebugLevel)
	tl, _ = log.NewLogger(cfg.LogPath+"traffic.log", zapcore.DebugLevel)
	sdkLogger, _ = log.NewLogger(cfg.LogPath+"sdk.log", zapcore.DebugLevel)
	tjLogger, _ = log.NewLogger(cfg.LogPath+"trojan.log", zapcore.DebugLevel)
}

func InitWebApi() {
	logger.Info("init mu api")
	cfg := cfg.WebApi
	apiClient = musdk.NewClient(cfg.Url, cfg.Token, cfg.NodeId, musdk.TypeV2ray, sdkLogger)
	apiClient.SetLogger(sdkLogger)
	go apiClient.UpdateTrafficDaemon()
	return
}
