package service

import (
	"github.com/aglowhy/sign/internal/app/bll/guilinlife"
	"github.com/aglowhy/sign/internal/app/config"
	"github.com/aglowhy/sign/pkg/logger"
)

func Service() {
	cfg := config.Global().GuilinlifeConf

	err := guilinlife.Exec(&cfg)
	if err != nil {
		logger.Errorf(nil, "桂林人论坛签到失败: %s", err.Error())
	}
}
