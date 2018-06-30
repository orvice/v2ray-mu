package main

import (
	"time"

	"github.com/catpie/musdk-go"
	"github.com/orvice/v2ray-manager"
)

func getV2rayManager() (*v2raymanager.Manager, error) {
	vm, err := v2raymanager.NewManager(cfg.V2rayClientAddr, cfg.V2rayTag)
	return vm, err
}

func (u *UserManager) check() error {
	logger.Info("Checking users from mu...")
	users, err := apiClient.GetUsers()
	if err != nil {
		logger.Errorf("Get users from error: %v", err)
		return err
	}
	logger.Infof("Get %d users from mu", len(users))
	for _, user := range users {
		u.checkUser(user)
	}

	return nil
}

func (u *UserManager) checkUser(user musdk.User) error {
	var err error
	if user.IsEnable() && !u.Exist(user) {
		logger.Infof("Run user id %d uuid %s", user.Id, user.V2rayUser.UUID)
		// run user
		err = u.vm.AddUser(&user.V2rayUser)
		if err != nil {
			logger.Errorf("Add user %s error %v", user.V2rayUser.UUID, err)
			return err
		}
		logger.Infof("Add user success %s", user.V2rayUser.UUID)
		u.AddUser(user)
		return nil
	}

	if user.IsEnable() && u.Exist(user) {
		ExistUser,ok:= u.GetUser(user.Id)
		if ExistUser.V2rayUser.UUID != user.V2rayUser.UUID && ok{
			logger.Infof("Changing user %d 's uuid to %s", user.Id, user.V2rayUser.UUID)
			err = u.vm.RemoveUser(&user.V2rayUser)

			if err != nil {
				logger.Errorf("Change user's UUID error %v", err)
				time.Sleep(time.Second * 10)
				return err
			}
			u.RemoveUser(user)
			err = u.vm.AddUser(&user.V2rayUser)
			if err != nil {
				logger.Errorf("Change user's UUID %s error %v", user.V2rayUser.UUID, err)
				return err
			}
			logger.Infof("Change user's UUID success %s", user.V2rayUser.UUID)
			u.AddUser(user)
		}
		return nil
	}

	if !user.IsEnable() && u.Exist(user) {
		logger.Infof("Stop user id %d uuid %s", user.Id, user.V2rayUser.UUID)
		// stop user
		err = u.vm.RemoveUser(&user.V2rayUser)

		if err != nil {
			logger.Errorf("Remove user error %v", err)
			time.Sleep(time.Second * 10)
			return err
		}
		u.RemoveUser(user)
		return nil
	}

	return nil
}

func (u *UserManager) restartUser() {}

func (u *UserManager) Run() error {
	for {
		u.saveTrafficDaemon()
		postNodeInfo()
		u.check()
		time.Sleep(cfg.SyncTime)
	}
	return nil
}

func (u *UserManager) Down() {
	u.cancel()
}

func (u *UserManager) saveTrafficDaemon() {
	logger.Infof("Runing save traffic daemon...")
	u.usersMu.RLock()
	defer u.usersMu.RUnlock()
	for _, user := range u.users {
		u.saveUserTraffic(user)
	}
}

func postNodeInfo() error {
	logger.Infof("Posting node info...")
	err := apiClient.PostNodeInfo()
	if err != nil {
		logger.Errorf("Post node info error %v", err)
	}
	return nil
}
