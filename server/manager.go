package server

import (
	"context"
	"strings"

	"github.com/catpie/musdk-go"
	v2raymanager "github.com/orvice/v2ray-manager"
	"github.com/p4gefau1t/trojan-go/api/service"
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

func getTrojanMgrs() ([]*TrojanMgr, error) {
	arr := strings.Split(cfg.TrojanApiServerAddr, ",")
	var trojanMgrs = make([]*TrojanMgr, len(arr))
	for k, v := range arr {
		tm, err := newTrojanMgr(v)
		if err != nil {
			return nil, err
		}
		trojanMgrs[k] = tm
	}
	return trojanMgrs, nil
}

func (u *UserManager) v2rayCheck() error {
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
		logger.Debugf("api users map add %s", user.V2rayUser.UUID)
		apiUsersMap[user.V2rayUser.UUID] = user
	}
	for _, user := range v2Users {
		v2UsersMap[user.User.GetUUID()] = user.User
	}

	// check remove user
	for _, v := range v2Users {
		uu, ok := apiUsersMap[v.User.GetUUID()]
		if !ok {
			logger.Infof("v2 user not in api, %s should be removed", v.User.GetUUID())
			u.vm.RemoveUser(ctx, v.User)
		}

		if uu.Enable == 0 {
			logger.Infof("user %s is disable, should be removed", v.User.GetUUID())
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
			logger.Infof("user %s may be should be add", v.V2rayUser.UUID)
			u.addUser(ctx, &v.V2rayUser)
		}
	}

	var logCount int

	logger.Infof("start v2 user data check len %d", len(v2Users))

	for _, vv := range v2Users {

		apiU, ok := apiUsersMap[vv.User.GetUUID()]
		if !ok {
			logger.Infof("%s is not found in api users ", vv.User.GetUUID())
			continue
		}

		if vv.TrafficInfo.Up == 0 && vv.TrafficInfo.Down == 0 {
			continue
		}

		logCount++

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

	logger.Infof("finish traffic log post len %d", logCount)

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
	switch u.targetType {
	case v2ray:
		task.NewTaskAndRun("check_users", cfg.SyncTime, u.v2rayCheck, task.SetTaskLogger(sdkLogger))
	case trojan:
		task.NewTaskAndRun("check_users", cfg.SyncTime, u.trojanCheck, task.SetTaskLogger(sdkLogger))
	}
	return nil

}

func (u *UserManager) Down() {
	u.cancel()
}

func (u *UserManager) trojanCheck() error {
	// ctx := context.Background()
	logger.Info("check users from mu")
	users, err := apiClient.GetUsers()
	if err != nil {
		tjLogger.Errorw("get users fail ",
			"error", err,
		)
		return err
	}
	tjLogger.Infof("get %d users from mu", len(users))

	// list users
	tus, err := u.tm.ListUsers()
	tjLogger.Infof("get %d users from trojan", len(tus))

	// add all users
	for _, user := range users {
		if user.Enable == 0 {
			continue
		}
		tjLogger.Infof("add trojan user %s", user.V2rayUser.UUID)
		err = u.tm.setUserStream.Send(&service.SetUsersRequest{
			Operation: service.SetUsersRequest_Add,
			Status: &service.UserStatus{
				User: &service.User{
					Password: user.V2rayUser.UUID,
				},
			},
		})
		if err != nil {
			tjLogger.Errorw("add trojan user error",
				"error", err,
			)
		}
		reply, err := u.tm.setUserStream.Recv()
		if err != nil {
			tjLogger.Errorw("fail to recv from set user stream",
				"error", err,
			)
		}

		tjLogger.Infof("add trojan user %s reply %v", user.V2rayUser.UUID, reply)

	}

	return nil
}
