package server

import (
	"github.com/catpie/musdk-go"
	"github.com/weeon/log"
	"golang.org/x/exp/slog"
)

var (
	apiClient *musdk.Client
)

var (
	logger    *slog.Logger
	tl        *slog.Logger // traffic logger
	sdkLogger *slog.Logger
	tjLogger  *slog.Logger
)

func Init() {
	initCfg()

	logger = slog.Default()
	tl = slog.Default()
	sdkLogger = slog.Default()
	tjLogger = slog.Default()
}

func InitWebApi() {
	logger.Info("init mu api")
	log.SetupStdoutLogger()
	cfg := cfg.WebApi
	apiClient = musdk.NewClient(cfg.Url, cfg.Token, cfg.NodeID, musdk.TypeV2ray, log.GetDefault())
	apiClient.SetLogger(log.GetDefault())
	go apiClient.UpdateTrafficDaemon()
	return
}
