package main

import (
	"github.com/catpie/musdk-go"
)

var (
	apiClient *musdk.Client
)

func InitWebApi() error {
	logger.Info("Initializing mu api")
	cfg := cfg.WebApi
	apiClient = musdk.NewClient(cfg.Url, cfg.Token, cfg.NodeId, musdk.TypeV2ray)
	apiClient.SetLogger(logger)
	go apiClient.UpdateTrafficDaemon()
	return nil
}
