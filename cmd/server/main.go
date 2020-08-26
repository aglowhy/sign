package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/aglowhy/sign/internal/app"
	"github.com/aglowhy/sign/pkg/logger"
	"github.com/aglowhy/sign/pkg/util"
)

// VERSION 版本号，
// 可以通过编译的方式指定版本号：go build -ldflags "-X main.VERSION=x.x.x"
var VERSION = "0.0.1"

var (
	configFile string
)

func init() {
	flag.StringVar(&configFile, "c", "", "配置文件(.toml)")
}

func main() {
	flag.Parse()

	if configFile == "" {
		panic("请使用-c指定配置文件")
	}

	var state int32 = 1
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	// 初始化日志参数
	logger.SetVersion(VERSION)
	logger.SetTraceIDFunc(util.NewTraceID)
	ctx := logger.NewTraceIDContext(context.Background(), util.NewTraceID())
	span := logger.StartSpanWithCall(ctx)

	call := app.Init(ctx,
		app.SetConfigFile(configFile),
		app.SetVersion(VERSION),
	)

EXIT:
	for {
		sig := <-sc
		span().Printf("获取到信号[%s]", sig.String())

		switch sig {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			atomic.StoreInt32(&state, 0)
			break EXIT
		case syscall.SIGHUP:
		default:
			break EXIT
		}
	}

	span().Printf("服务退出")

	if call != nil {
		call()
	}

	time.Sleep(time.Second)
	os.Exit(int(atomic.LoadInt32(&state)))
}
