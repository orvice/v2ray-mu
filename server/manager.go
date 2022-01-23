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

		if vv.TrafficInfo.Up == 0 && vv.TrafficInfo.Down == 0 {
			continue
		}

		logCount++

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
	ctx, cancel := context.WithCancel(u.ctx)
	defer cancel()

	logger.Info("[trojan] check users from mu")
	users, err := apiClient.GetUsers()
	if err != nil {
		tjLogger.Errorw("[trojan] get users fail ",
			"error", err,
		)
		return err
	}
	tjLogger.Infof("[trojan] get %d users from mu", len(users))

	stream, err := u.tm.client.SetUsers(ctx)
	if err != nil {
		return err
	}

	getUserClient, err := u.tm.client.GetUsers(ctx)
	if err != nil {
		return err
	}

	var trafficLogCount int

	// add all users
	for _, user := range users {
		resp, err := u.tm.GetUser(ctx, getUserClient, user.V2rayUser.UUID)
		tjLogger.Infof("[trojan] get user reploy %v", resp)

		if resp != nil && resp.Status != nil {
			status, ok := u.tm.userStatusMap[user.Id]
			if ok {
				u := resp.Status.TrafficTotal.UploadTraffic - status.TrafficTotal.UploadTraffic
				d := resp.Status.TrafficTotal.DownloadTraffic - status.TrafficTotal.DownloadTraffic
				if u > 0 && d > 0 {
					trafficLog := musdk.UserTrafficLog{
						UserId: user.Id,
						U:      int64(u),
						D:      int64(d),
					}

					tl.Infow("[trojan] save raffice log",
						"user_id", user.Id,
						"traffic_log", trafficLog,
					)
					trafficLogCount++
					apiClient.SaveTrafficLog(trafficLog)
				}
			}

			u.tm.userStatusMap[user.Id] = resp.Status
		}

		if user.Enable == 0 {

			if err != nil {
				continue
			}

			if resp == nil {
				continue
			}

			// remove user
			err = stream.Send(&service.SetUsersRequest{
				Operation: service.SetUsersRequest_Delete,
				Status: &service.UserStatus{
					User: &service.User{
						Password: user.V2rayUser.UUID,
					},
				},
			})
			if err != nil {
				tjLogger.Errorf("[trojan] trojan remove user %s error %v", user.V2rayUser.UUID, err)
			}
			reply, err := stream.Recv()
			if err != nil {
				tjLogger.Errorw("[trojan] fail to recv from set user stream",
					"error", err,
				)
			}

			tjLogger.Infof("delete trojan user %s reply %v", user.V2rayUser.UUID, reply)

			continue
		}

		if err == nil && resp.Status.User != nil {
			tjLogger.Infof("[trojan] user %s exist", user.V2rayUser.UUID)
			continue
		}

		tjLogger.Infof("[trojan] add trojan user %s", user.V2rayUser.UUID)
		err = stream.Send(&service.SetUsersRequest{
			Operation: service.SetUsersRequest_Add,
			Status: &service.UserStatus{
				User: &service.User{
					Password: user.V2rayUser.UUID,
				},
			},
		})
		if err != nil {
			tjLogger.Errorw("[trojan] add trojan user error",
				"error", err,
			)
		}
		reply, err := stream.Recv()
		if err != nil {
			tjLogger.Errorw("[trojan] fail to recv from set user stream",
				"error", err,
			)
		}

		tjLogger.Infof("[trojan] add trojan user %s reply %v", user.V2rayUser.UUID, reply)

	}

	tjLogger.Infof("traffic log count %d", trafficLogCount)

	return nil
}
