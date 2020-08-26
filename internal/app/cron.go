package app

import (
	"context"
	"github.com/aglowhy/sign/internal/app/service"
	"github.com/aglowhy/sign/pkg/logger"
	"github.com/aglowhy/sign/pkg/util"
	"github.com/robfig/cron/v3"
	"time"
)

// InitCron 初始化crontab服务
func InitCron(ctx context.Context) func() {
	c := cron.New(cron.WithSeconds())

	go func() {
		logger.Printf(ctx, "Cron服务开始启动...")

		spec := "10 0 0 * * *"
		_, err := c.AddFunc(spec, func() {
			r := util.RandomInt(10, 300)
			logger.Printf(ctx, "任务将会在 %d s后执行", r)
			time.Sleep(time.Duration(r) * time.Second)

			service.Service()
		})
		if err != nil {
			logger.Errorf(ctx, err.Error())
		}

		c.Start()
	}()

	return func() {
		c.Stop()
	}
}
