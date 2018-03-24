package main

import (
	"time"

	"github.com/orvice/v2ray-manager"
)

func getV2rayManager() (*v2raymanager.Manager, error) {
	vm, err := v2raymanager.NewManager(cfg.V2rayClientAddr, cfg.V2rayTag)
	return vm, err
}

func (u *UserManager) check() error {
	logger.Info("check users from mu")
	users, err := apiClient.GetUsers()
	if err != nil {
		logger.Errorf("get users from error: %v", err)
		return err
	}
	logger.Infof("get %d users from mu", len(users))
	for _, user := range users {
		if user.IsEnable() && !u.Exist(user) {
			logger.Infof("run user id %d uuid %s", user.Id, user.V2rayUser.UUID)
			// run user
			err = u.vm.AddUser(&user.V2rayUser)
			if err != nil {
				logger.Errorf("add user error %v", err)
				// @todo error handle
				time.Sleep(time.Second * 10)
				continue
			}
			u.AddUser(user)
			continue
		}

		if !user.IsEnable() && u.Exist(user) {
			logger.Infof("stop user id %d uuid %s", user.Id, user.V2rayUser.UUID)
			// stop user
			err = u.vm.RemoveUser(&user.V2rayUser)

			if err != nil {
				logger.Errorf("remove user error %v", err)
				time.Sleep(time.Second * 10)
				continue
				// @todo error handle
			}
			u.RemoveUser(user)
			continue
		}
	}

	return nil
}

func (u *UserManager) restartUser() {}

func (u *UserManager) Run() error {
	for {
		u.check()
		time.Sleep(cfg.SyncTime)
	}
	return nil
}

func (u *UserManager) Down() {
	u.cancel()
}
