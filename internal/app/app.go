package app

import (
	"context"
	"github.com/aglowhy/sign/internal/app/config"
	"github.com/aglowhy/sign/pkg/logger"
	"os"
)

type options struct {
	ConfigFile string
	Version    string
}

// Option 定义配置项
type Option func(*options)

// SetConfigFile 设定配置文件
func SetConfigFile(s string) Option {
	return func(o *options) {
		o.ConfigFile = s
	}
}

// SetVersion 设定版本号
func SetVersion(s string) Option {
	return func(o *options) {
		o.Version = s
	}
}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}

func Init(ctx context.Context, opts ...Option) func() {
	var o options
	for _, opt := range opts {
		opt(&o)
	}
	err := config.LoadGlobal(o.ConfigFile)
	handleError(err)

	cfg := config.Global()

	loggerCall, err := InitLogger()
	handleError(err)

	logger.Printf(ctx, "服务启动，运行模式：%s，版本号：%s，进程号：%d", cfg.RunMode, o.Version, os.Getpid())

	cronCall := InitCron(ctx)

	return func() {
		if loggerCall != nil {
			loggerCall()
		}
		if cronCall != nil {
			cronCall()
		}
	}
}
