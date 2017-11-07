package main

import (
	"github.com/catpie/musdk-go"
	"github.com/orvice/utils/log"
)

var (
	apiClient *musdk.Client
)

func InitWebApi() error {
	log.Info("init mu api")
	cfg := cfg.WebApi
	apiClient = musdk.NewClient(cfg.Url, cfg.Token, cfg.NodeId, musdk.TypeSs)
	return nil
}
