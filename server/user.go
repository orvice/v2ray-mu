package server

import (
	"context"
	"sync"

	"github.com/catpie/musdk-go"
	v2raymanager "github.com/orvice/v2ray-manager"
)

const (
	v2ray  = 0
	trojan = 1
)

type UserManager struct {
	users   map[int64]musdk.User
	usersMu *sync.RWMutex
	ctx     context.Context
	cancel  func()

	vm         *v2raymanager.Manager
	tm         *TrojanMgr
	targetType int
}

func NewUserManager() ([]*UserManager, error) {
	ctx, cancel := context.WithCancel(context.Background())
	vms, err := getV2rayManager()
	if err != nil {
		return nil, err
	}
	tms, err := getTrojanMgrs()
	if err != nil {
		return nil, err
	}

	ums := make([]*UserManager, len(vms)+len(tms))
	for k, v := range vms {
		ums[k] = &UserManager{
			users:      make(map[int64]musdk.User),
			usersMu:    new(sync.RWMutex),
			ctx:        ctx,
			cancel:     cancel,
			vm:         v,
			targetType: v2ray,
		}
	}

	for k, v := range tms {
		ums[len(vms)+k] = &UserManager{
			users:      make(map[int64]musdk.User),
			usersMu:    new(sync.RWMutex),
			ctx:        ctx,
			cancel:     cancel,
			tm:         v,
			targetType: trojan,
		}
	}

	return ums, nil
}

func (u *UserManager) GetUser(id int64) (musdk.User, bool) {
	user, ok := u.users[id]
	return user, ok
}
