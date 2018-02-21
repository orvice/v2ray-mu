package main

import (
	"time"

	"github.com/orvice/v2ray-manager"
	log "github.com/sirupsen/logrus"
)

var (
	VM *v2raymanager.Manager
)

func InitV2rayManager() error {
	var err error
	VM, err = v2raymanager.NewManager(cfg.V2rayClientAddr, cfg.V2rayTag)
	return err
}

func check() error {
	log.Info("check users from mu")
	users, err := apiClient.GetUsers()
	if err != nil {
		return err
	}
	log.Infof("get %d users from mu", len(users))
	for _, user := range users {
		if user.IsEnable() && !UM.Exist(user) {
			log.Infof("run user id %d uuid %s", user.Id, user.V2rayUser.UUID)
			// run user
			err = VM.AddUser(&user.V2rayUser)
			if err != nil {
				// @todo error handle
			}
			UM.AddUser(user)
			continue
		}

		if !user.IsEnable() && UM.Exist(user) {
			log.Infof("stop user id %d uuid %s", user.Id, user.V2rayUser.UUID)
			// stop user
			err = VM.RemoveUser(&user.V2rayUser)

			if err != nil {
				// @todo error handle
			}
			UM.RemoveUser(user)
			continue
		}
	}

	return nil
}

func daemon() {
	for {
		check()
		time.Sleep(cfg.SyncTime)
	}
}
