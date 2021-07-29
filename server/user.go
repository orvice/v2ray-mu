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

func NewUserManager() ([]*UserManager, error) {
	ctx, cancel := context.WithCancel(context.Background())
	vms, err := getV2rayManager()
	if err != nil {
		return nil, err
	}
	ums := make([]*UserManager, len(vms))
	for k, v := range vms {
		ums[k] = &UserManager{
			users:   make(map[int64]musdk.User),
			usersMu: new(sync.RWMutex),
			ctx:     ctx,
			cancel:  cancel,
			vm:      v,
		}
	}
	return ums, nil
}

func (u *UserManager) GetUser(id int64) (musdk.User, bool) {
	user, ok := u.users[id]
	return user, ok
}
