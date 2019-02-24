package server

import (
	"context"
	"sync"

	"github.com/catpie/musdk-go"
	"github.com/orvice/v2ray-manager"
)

type UserManager struct {
	users   map[int64]musdk.User
	usersMu *sync.RWMutex
	ctx     context.Context
	cancel  func()

	vm *v2raymanager.Manager
}

func NewUserManager() (*UserManager, error) {
	ctx, cancel := context.WithCancel(context.Background())
	vm, err := getV2rayManager()
	if err != nil {
		return nil, err
	}
	um := &UserManager{
		users:   make(map[int64]musdk.User),
		usersMu: new(sync.RWMutex),
		ctx:     ctx,
		cancel:  cancel,
		vm:      vm,
	}
	return um, nil
}

func (u *UserManager) AddUser(user musdk.User) {
	u.usersMu.Lock()
	u.users[user.Id] = user
	u.usersMu.Unlock()
}

func (u *UserManager) RemoveUser(user musdk.User) {
	u.usersMu.Lock()
	delete(u.users, user.Id)
	u.usersMu.Unlock()
}

func (u *UserManager) GetUser(id int64) (musdk.User, bool) {
	user, ok := u.users[id]
	return user, ok
}

func (u *UserManager) Exist(user musdk.User) bool {
	u.usersMu.RLock()
	defer u.usersMu.RUnlock()
	_, ok := u.users[user.Id]
	if ok {
		return true
	}
	return false
}

func (u *UserManager) saveUserTraffic(user musdk.User) {
	ti := u.vm.GetTrafficAndReset(&user.V2rayUser)
	if ti.Down == 0 && ti.Up == 0 {
		return
	}
	trafficLog := musdk.UserTrafficLog{
		UserId: user.Id,
		U:      ti.Up,
		D:      ti.Down,
	}
	tl.Infof("id %d uuid %s save traffice log %v", user.Id, user.V2rayUser.UUID, trafficLog)
	apiClient.SaveTrafficLog(trafficLog)
}
