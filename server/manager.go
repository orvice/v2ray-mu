package server

import (
	"github.com/catpie/musdk-go"
	"github.com/orvice/v2ray-manager"
)

func getV2rayManager() (*v2raymanager.Manager, error) {
	vm, err := v2raymanager.NewManager(cfg.V2rayClientAddr, cfg.V2rayTag)
	vm.SetLogger(logger)
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
	logger.Infof("reload user")
	for _, user := range users {
		u.checkUser(user)
	}
	logger.Infof("reload user finished")
	return nil
}

func (u *UserManager) checkUser(user musdk.User) error  {
		// run user
	logger.Infof("stop user id %d uuid %s", user.Id, user.V2rayUser.UUID)
		// stop user
		err = u.vm.RemoveUser(&user.V2rayUser)
		if err != nil {
			logger.Errorf("remove user error %v", err)
			return err
		}
		u.RemoveUser(user)
	exist, err := u.vm.AddUser(&user.V2rayUser)
		if err != nil {
			logger.Errorf("add user %s error %v", user.V2rayUser.UUID, err)
			return err
		}
		if !exist {
			logger.Errorf("add user %s success", user.V2rayUser.UUID)
		}
		u.AddUser(user)
	return nil
}

func (u *UserManager) restartUser() {}

func (u *UserManager) Run() error {
	runJob("check_users", cfg.SyncTime, u.check)
	runJob("save_traffic", cfg.SyncTime, u.saveTrafficDaemon)
	return nil

}

func (u *UserManager) Down() {
	u.cancel()
}

func (u *UserManager) saveTrafficDaemon() error {
	u.usersMu.RLock()
	defer u.usersMu.RUnlock()
	for _, user := range u.users {
		u.saveUserTraffic(user)
	}
	return nil
}
