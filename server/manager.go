package server

import (
	"time"

	"github.com/catpie/musdk-go"
	v2raymanager "github.com/orvice/v2ray-manager"
	"github.com/weeon/utils/task"
)

func getV2rayManager() (*v2raymanager.Manager, error) {
	vm, err := v2raymanager.NewManager(cfg.V2rayClientAddr, cfg.V2rayTag, sdkLogger)
	return vm, err
}

func (u *UserManager) check() error {
	logger.Info("check users from mu")
	users, err := apiClient.GetUsers()
	if err != nil {
		logger.Errorw("get users fail ",
			"error", err,
		)
		return err
	}
	logger.Infof("get %d users from mu", len(users))

	apiUserMap := make(map[int64]bool)
	// update users
	for _, user := range users {
		err := u.checkUser(user)
		if err != nil {
			logger.Errorf("check user(id=%d) fail", user.GetId())
		}
		apiUserMap[user.GetId()] = true
	}

	// sync users
	if len(users) != len(u.GetUsers()) {
		for uid, user := range u.GetUsers() {
			_, ok := apiUserMap[uid]
			if !ok {
				logger.Infof("user(id=%d) no longer exists in mu system, removing it", uid)
				// user no longer exists in mu system, remove it
				err = u.vm.RemoveUser(&user.V2rayUser)
				if err != nil {
					logger.Errorf("remove user error %v", err)
					time.Sleep(time.Second * 10)
					return err
				}
				u.RemoveUser(user)
			}
		}
	}

	return nil
}

func (u *UserManager) checkUser(user musdk.User) error {
	var err error
	if user.IsEnable() {
		if !u.Exist(user) {
			logger.Infof("add user id %d uuid %s", user.Id, user.V2rayUser.UUID)
			// add user
			exist, err := u.vm.AddUser(&user.V2rayUser)
			if err != nil {
				logger.Errorf("add user %s error %v", user.V2rayUser.UUID, err)
				return err
			}
			if exist {
				logger.Infof("user %s already exist", user.V2rayUser.UUID)
			}
			u.AddUser(user)
			return nil
		} else {
			// update user info
			oldUser, ok := u.GetUser(user.GetId())
			if ok {
				// check if user info needs update
				if oldUser.V2rayUser.GetAlterID() != user.V2rayUser.GetAlterID() ||
					oldUser.V2rayUser.GetEmail() != user.V2rayUser.GetEmail() ||
					oldUser.V2rayUser.GetLevel() != user.V2rayUser.GetLevel() ||
					oldUser.V2rayUser.GetUUID() != user.V2rayUser.GetUUID() {
					logger.Infof("update user id %d old uuid %s uuid %s",
						user.GetId(), oldUser.V2rayUser.GetUUID(), user.V2rayUser.GetUUID())
					// do update, first remove then add
					// remove old user
					err = u.vm.RemoveUser(&oldUser.V2rayUser)
					if err != nil {
						logger.Errorf("remove user error %v", err)
						time.Sleep(time.Second * 10)
						return err
					}
					u.RemoveUser(oldUser)

					// add new user
					_, err := u.vm.AddUser(&user.V2rayUser)
					if err != nil {
						logger.Errorf("add user %s error %v", user.V2rayUser.UUID, err)
						return err
					}
					u.AddUser(user)
					return nil
				}
			}
		}
	}

	if !user.IsEnable() && u.Exist(user) {
		logger.Infof("stop user id %d uuid %s", user.Id, user.V2rayUser.UUID)
		// stop user
		err = u.vm.RemoveUser(&user.V2rayUser)
		if err != nil {
			logger.Errorf("remove user error %v", err)
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
	task.NewTaskAndRun("check_users", cfg.SyncTime, u.check, task.SetTaskLogger(sdkLogger))
	task.NewTaskAndRun("save_traffic", cfg.SyncTime, u.saveTrafficDaemon, task.SetTaskLogger(sdkLogger))
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
