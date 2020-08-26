package notice

import (
	"context"
	"net/http"
	"net/url"
	"strings"

	"github.com/aglowhy/sign/internal/app/config"
	"github.com/aglowhy/sign/pkg/email"
	"github.com/aglowhy/sign/pkg/logger"
)

type Message struct {
	Title string
	Desc  string
}

// 消息通知
func ServerChanNotice(msg *Message) {
	data := url.Values{
		"text": {msg.Title},
		"desp": {msg.Desc},
	}
	u := config.Global().ServerChan.Url
	body := strings.NewReader(data.Encode())
	resp, err := http.Post(u, "application/x-www-form-urlencoded", body)
	if err != nil {
		logger.Errorf(context.Background(), "Server Chan notice fail: %s", err.Error())
		msg := &email.Message{
			Subject: "Server酱消息推送失败",
			Body:    "Error msg: " + err.Error(),
		}
		email.SendEmail(msg)
		panic(err)
	}

	defer resp.Body.Close()
}
