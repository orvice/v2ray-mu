package main

import (
	"github.com/orvice/v2ray-manager"
	"time"
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
	users, err := apiClient.GetUsers()
	if err != nil {
		return err
	}
	for _, user := range users {
		if user.IsEnable() && !UM.Exist(user) {
			// run user
			err = VM.AddUser(&user.V2rayUser)
			if err != nil {
				// @todo error handle
			}
			UM.AddUser(user)
			continue
		}

		if !user.IsEnable() && UM.Exist(user) {
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
