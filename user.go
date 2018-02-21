package main

import (
	"sync"

	"github.com/catpie/musdk-go"
)

var (
	UM *UserManager
)

func InitUserManager() {
	UM = NewUserManager()
}

type UserManager struct {
	users   map[int64]musdk.User
	usersMu *sync.RWMutex
}

func NewUserManager() *UserManager {
	um := &UserManager{
		users:   make(map[int64]musdk.User),
		usersMu: new(sync.RWMutex),
	}
	return um
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

func (u *UserManager) Exist(user musdk.User) bool {
	u.usersMu.RLock()
	defer u.usersMu.RUnlock()
	_, ok := u.users[user.Id]
	if ok {
		return true
	}
	return false
}
