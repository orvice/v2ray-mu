package server

import (
	"context"
	"strings"

	"github.com/catpie/musdk-go"
	v2raymanager "github.com/orvice/v2ray-manager"
	"github.com/weeon/log"
	"github.com/weeon/utils/task"
)

func getV2rayManager() ([]*v2raymanager.Manager, error) {
	arr := strings.Split(cfg.V2rayClientAddr, ",")
	var vms = make([]*v2raymanager.Manager, len(arr))
	for k, v := range arr {
		vm, err := v2raymanager.NewManager(v, cfg.V2rayTag, log.GetDefault())
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
		logger.Error("get users fail ",
			"error", err,
		)
		return err
	}
	logger.Info("get  users from mu", "len", len(users))

	v2Users, err := u.vm.GetUserList(ctx, true)
	if err != nil {
		logger.Error("get users list error",
			"error", err.Error())
	}

	apiUsersMap := make(map[string]musdk.User)
	v2UsersMap := make(map[string]v2raymanager.User)

	for _, user := range users {
		logger.Debug("api users map add ", "uuid", user.V2rayUser.UUID)
		apiUsersMap[user.V2rayUser.UUID] = user
	}
	for _, user := range v2Users {
		v2UsersMap[user.User.GetUUID()] = user.User
	}

	// check remove user
	for _, v := range v2Users {
		uu, ok := apiUsersMap[v.User.GetUUID()]
		if !ok {
			logger.Info("v2 user not in api  should be removed", "uuid", v.User.GetUUID())
			u.vm.RemoveUser(ctx, v.User)
		}

		if uu.Enable == 0 {
			logger.Info("user  should be removed", "uuid", v.User.GetUUID())
			u.vm.RemoveUser(ctx, v.User)
		}

	}

	for _, v := range users {
		if v.Enable == 0 {
			continue
		}

		_, ok := v2UsersMap[v.V2rayUser.UUID]
		if !ok {
			logger.Info("user may be should be add", "uuid", v.V2rayUser.UUID)
			u.addUser(ctx, &v.V2rayUser)
		}
	}

	var logCount int

	logger.Info("start v2 user data check len ", "len", len(v2Users))

	for _, vv := range v2Users {

		apiU, ok := apiUsersMap[vv.User.GetUUID()]
		if !ok {
			logger.Info("not found in api users ", "uuid", vv.User.GetUUID())
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

		tl.Info("save traffice log",
			"user_id", apiU.Id,
			"uuid", vv.User.GetUUID(),
			"traffic Log", trafficLog,
		)
		apiClient.SaveTrafficLog(trafficLog)
	}

	logger.Info("finish traffic log post len", "len", logCount)

	return nil
}

func (u *UserManager) addUser(ctx context.Context, user v2raymanager.User) {
	exist, err := u.vm.AddUser(ctx, user)
	if err != nil {
		logger.Error("add user  error", "uuid", user.GetUUID(), "err", err)
		return
	}
	if !exist {
		logger.Info("add user  success", "uuid", user.GetUUID())
	}

	logger.Info("add user  result: AlreadyExists", "uuid", user.GetUUID())
	return
}

func (u *UserManager) Run() error {
	switch u.targetType {
	case v2ray:
		task.NewTaskAndRun("check_users", cfg.SyncTime, u.v2rayCheck, task.SetTaskLogger(log.GetDefault()))
	case trojan:
		task.NewTaskAndRun("check_users", cfg.SyncTime, u.trojanCheck, task.SetTaskLogger(log.GetDefault()))
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
		tjLogger.Error("[trojan] get users fail ",
			"error", err,
		)
		return err
	}
	tjLogger.Info("[trojan] get users from mu", "len", len(users))

	var trafficLogCount int
	var trafficLogs = make([]musdk.UserTrafficLog, 0)

	// add all users
	for _, user := range users {
		tjLogger.Info("[trojan] start get user ", "id", user.Id, "uuid", user.V2rayUser.UUID)
		resp, err := u.tm.GetUser(ctx, user.V2rayUser.UUID)
		tjLogger.Info(" [trojan] get user reploy ", "resp", resp, "uuid", user.V2rayUser.UUID)

		if err != nil {
			tjLogger.Error("[trojan] get user fail ",
				"error", err,
			)
			continue
		}

		if resp != nil && resp.Status != nil {
			status, ok := u.tm.userStatusMap[user.Id]
			if ok {
				u := int64(resp.Status.TrafficTotal.UploadTraffic) - int64(status.TrafficTotal.UploadTraffic)
				d := int64(resp.Status.TrafficTotal.DownloadTraffic) - int64(status.TrafficTotal.DownloadTraffic)
				if u > 0 && d > 0 {
					trafficLog := musdk.UserTrafficLog{
						UserId: user.Id,
						U:      int64(u),
						D:      int64(d),
					}

					trafficLogs = append(trafficLogs, trafficLog)

					tjLogger.Info("[trojan] save raffice log",
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

			tjLogger.Info("[trojan] user  is disable", "id", user.Id, "uuid", user.V2rayUser.UUID)

			if resp == nil {
				continue
			}

			if resp.Status == nil {
				tjLogger.Info("[trojan] user   status is nil", "id", user.Id, "uuid", user.V2rayUser.UUID)
				continue
			}

			// remove user
			logger.Info("[trojan] start remove user ", "id", user.Id, "uuid", user.V2rayUser.UUID, "user.hash", resp.Status.User.Hash)
			err = u.tm.RemoveUser(ctx, user.V2rayUser.UUID, resp.Status.User.Hash)
			if err != nil {
				tjLogger.Error("[trojan] trojan remove ", "uuid", user.V2rayUser.UUID, "err", err)
				continue
			}

			tjLogger.Info("delete trojan user  success ", "uuid", user.V2rayUser.UUID)

			continue
		}

		tjLogger.Info("[trojan] user  is enable", "id", user.Id, "uuid", user.V2rayUser.UUID)

		tjLogger.Info("check user is exist ", "uuid", user.V2rayUser.UUID)
		if resp.Success && resp.Status != nil {
			tjLogger.Info("[trojan] user  exist", "uuid", user.V2rayUser.UUID)
			continue
		}

		tjLogger.Info("[trojan] add trojan user ", "uuid", user.V2rayUser.UUID)
		err = u.tm.AddUser(ctx, user.V2rayUser.UUID)
		if err != nil {
			tjLogger.Error("[trojan] add trojan user error",
				"error", err,
			)
			continue
		}

		tjLogger.Info("[trojan] add trojan user  success", "uuid", user.V2rayUser.UUID)
	}

	tjLogger.Info("traffic log count ", "count", trafficLogCount)
	tjLogger.Info("traffic logs",
		"traffic_logs", trafficLogs,
	)

	return nil
}
