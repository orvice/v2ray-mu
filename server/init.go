package server

import (
	"os"

	"log/slog"

	"github.com/catpie/musdk-go"
	"github.com/weeon/log"
)

var (
	apiClient *musdk.Client
)

var (
	logger *slog.Logger
	tl     *slog.Logger // traffic logger
	// sdkLogger *slog.Logger
	tjLogger *slog.Logger
)

func Init() {
	initCfg()

	textHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{})
	logger = slog.New(textHandler)

	tl = slog.New(textHandler)
	// sdkLogger = slog.New(textHandler)
	tjLogger = slog.New(textHandler)
}

func InitWebApi() {
	logger.Info("init mu api")
	log.SetupStdoutLogger()
	cfg := cfg.WebApi
	apiClient = musdk.NewClient(cfg.Url, cfg.Token, cfg.NodeID, musdk.TypeV2ray, logger)
	apiClient.SetLogger(logger)
	go apiClient.UpdateTrafficDaemon()
}
