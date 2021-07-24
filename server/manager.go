package server

import (
	"context"
	"strings"
	"time"

	"github.com/catpie/musdk-go"
	v2raymanager "github.com/orvice/v2ray-manager"
	"github.com/weeon/utils/task"
)

func getV2rayManager() ([]*v2raymanager.Manager, error) {
	arr := strings.Split(cfg.V2rayClientAddr, ",")
	var vms = make([]*v2raymanager.Manager, len(arr))
	for k, v := range arr {
		vm, err := v2raymanager.NewManager(v, cfg.V2rayTag, sdkLogger)
		if err != nil {
			return nil, err
		}
		vms[k] = vm
	}

	return vms, nil
}

func (u *UserManager) check() error {
	ctx := context.Background()
	logger.Info("check users from mu")
	users, err := apiClient.GetUsers()
	if err != nil {
		logger.Errorw("get users fail ",
			"error", err,
		)
		return err
	}
	logger.Infof("get %d users from mu", len(users))
	for _, user := range users {
		u.checkUser(ctx, user)
	}

	return nil
}

func (u *UserManager) checkUser(ctx context.Context, user musdk.User) error {
	var err error
	if user.IsEnable() && !u.Exist(user) {
		// run user
		exist, err := u.vm.AddUser(ctx, &user.V2rayUser)
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

	if !user.IsEnable() && u.Exist(user) {
		logger.Infof("stop user id %d uuid %s", user.Id, user.V2rayUser.UUID)
		// stop user
		err = u.vm.RemoveUser(ctx, &user.V2rayUser)

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
	ctx := context.Background()
	u.usersMu.RLock()
	defer u.usersMu.RUnlock()
	for _, user := range u.users {
		u.saveUserTraffic(ctx, user)
	}
	return nil
}
