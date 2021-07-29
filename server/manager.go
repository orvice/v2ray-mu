package server

import (
	"context"
	"github.com/catpie/musdk-go"
	v2raymanager "github.com/orvice/v2ray-manager"
	"github.com/weeon/log"
	"github.com/weeon/utils/task"
	"strings"
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

	v2Users, err := u.vm.GetUserList(ctx, true)
	if err != nil {
		logger.Errorw("get users list error",
			"error", err.Error())
	}

	apiUsersMap := make(map[string]musdk.User)
	v2UsersMap := make(map[string]v2raymanager.User)

	for _, user := range users {
		apiUsersMap[user.V2rayUser.GetUUID()] = user
	}
	for _, user := range v2Users {
		v2UsersMap[user.User.GetUUID()] = user.User
	}

	// check remove user
	for _, v := range v2Users {
		uu, ok := apiUsersMap[v.User.GetUUID()]
		if !ok {
			log.Infof("v2 user not in api, %s should be removed", v.User.GetUUID())
			u.vm.RemoveUser(ctx, v.User)
		}

		if uu.Enable == 0 {
			log.Infof("user %s is disable, should be removed", v.User.GetUUID())
			u.vm.RemoveUser(ctx, v.User)
		}

	}

	// check add
	for _, v := range users {
		if v.Enable == 0 {
			continue
		}

		_, ok := v2UsersMap[v.V2rayUser.UUID]
		if !ok {
			log.Infof("user %s may be should be add", v.V2rayUser.UUID)
			u.addUser(ctx, &v.V2rayUser)
		}
	}

	for _, vv := range v2Users {

		apiU, ok := apiUsersMap[vv.User.GetUUID()]
		if ok {
			continue
		}

		trafficLog := musdk.UserTrafficLog{
			UserId: apiU.Id,
			U:      vv.TrafficInfo.Up,
			D:      vv.TrafficInfo.Down,
		}
		tl.Infow("save traffice log",
			"user_id", apiU.Id,
			"uuid", vv.User.GetUUID(),
			"traffic Log", trafficLog,
		)
		apiClient.SaveTrafficLog(trafficLog)
	}

	return nil
}

func (u *UserManager) addUser(ctx context.Context, user v2raymanager.User) {
	exist, err := u.vm.AddUser(ctx, user)
	if err != nil {
		logger.Errorf("add user %s error %v", user.GetUUID(), err)
		return
	}
	if !exist {
		logger.Errorf("add user %s success", user.GetUUID())
	}
	return
}

func (u *UserManager) Run() error {
	task.NewTaskAndRun("check_users", cfg.SyncTime, u.check, task.SetTaskLogger(sdkLogger))
	return nil

}

func (u *UserManager) Down() {
	u.cancel()
}
