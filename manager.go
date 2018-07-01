package main

import (
	"time"
	"errors"

	"github.com/catpie/musdk-go"
	"github.com/orvice/v2ray-manager"
	"github.com/orvice/shadowsocks-go/mu/system"
	"net/http"
	"fmt"
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
		u.postNodeInfo()
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

func (u *UserManager) postNodeInfo() error {
	logger.Infof("Posting node info...")
	err := u.PostNodeInfo()
	if err != nil {
		logger.Errorf("Post node info error %v", err)
	}
	return nil
}

func (u *UserManager) postNodeInfoUri() string {
	return fmt.Sprintf("%s/nodes/%d/info", cfg.WebApi.Url, cfg.WebApi.NodeId)
}

func (u *UserManager) PostNodeInfo() error {
	uptime, err := system.GetUptime()
	if err != nil {
		uptime = "0"
	}

	load, err := system.GetLoad()
	if err != nil {
		load = "- - -"
	} else {
	load = load[0:13]
	}
	data := `{"load":"`+load+`","uptime":"`+uptime+`"}`
	_, statusCode, err := u.httpPost(u.postNodeInfoUri(), string(data))
	if err != nil {
		return err
	}
	if statusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("status code: %d", statusCode))
	}
	return nil
}
