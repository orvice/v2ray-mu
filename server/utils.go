package server

import (
	"strings"
	"time"

	"github.com/orvice/kit/utils"
)

func IsAlreadyExistsError(e error) bool {
	if strings.Contains(e.Error(), "already exists") {
		return true
	}
	return false
}

func runJob(name string, t time.Duration, fn func() error) {
	task := utils.NewTask(name, t, fn, func(opt *utils.Task) {
		opt.Logger = logger
	})
	go task.Run()
}
