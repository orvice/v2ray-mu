package server

import (
	"github.com/shuangzhijinghua/musdk-go"
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
	for _, user := range users {
		u.checkUser(user)
	}
	return nil
}

func (u *UserManager) checkUser(user musdk.User) error  {
	var traffic, maxtraffic int64
	maxtraffic = int64(cfg.MaxTraffic) * 1024 * 1024 * 1024
	traffic = user.U + user.D
	if ( traffic >= maxtraffic ) && u.Exist(user){
		u.vm.RemoveUser(&user.V2rayUser)
		u.RemoveUser(user)
	} else if !u.UserDiff(user) {
		logger.Infof("user id %d email %s uuid %s  is different with previous one,reloading.", user.Id, user.V2rayUser.Email, user.V2rayUser.UUID)
		u.vm.RemoveUser(&user.V2rayUser)
		u.RemoveUser(user)
	}

	if ( traffic < maxtraffic ) && user.IsEnable() && !u.Exist(user) {
		logger.Infof("user %s is valid, current %v GiB, will be add to v2ray.", user.V2rayUser.Email, int(traffic/1024/1024/1024))
		u.vm.AddUser(&user.V2rayUser)
		u.AddUser(user)
	} else if user.Admin() {
		logger.Infof("user %d is admin, add anyway.",user.Id)
		u.vm.AddUser(&user.V2rayUser)
		u.AddUser(user)
	} else if ( traffic >= maxtraffic) {
		logger.Infof("user %s is overusage, current %v GiB, will not add to v2ray.", user.V2rayUser.Email, int(traffic/1024/1024/1024))
	}
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
