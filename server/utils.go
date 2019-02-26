package server

import (
	"time"

	"github.com/orvice/kit/utils"
)


func runJob(name string, t time.Duration, fn func() error) {
	task := utils.NewTask(name, t, fn, func(opt *utils.Task) {
		opt.Logger = logger
	})
	go task.Run()
}
