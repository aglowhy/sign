package email

import (
	"context"
	"github.com/aglowhy/sign/internal/app/config"
	"github.com/aglowhy/sign/pkg/logger"
	"github.com/go-mail/mail"
)

type Message struct {
	Subject string
	Body    string
}

// 邮件通知
func SendEmail(msg *Message) {
	c := config.Global().Email

	m := mail.NewMessage()
	m.SetHeader("From", c.Host)
	m.SetHeader("To", c.Receiver)
	m.SetHeader("Subject", msg.Subject)
	m.SetBody("text/html", msg.Body)

	d := mail.NewDialer(c.Host, c.Port, c.UserName, c.Password)
	//d.StartTLSPolicy = mail.MandatoryStartTLS

	if err := d.DialAndSend(m); err != nil {
		logger.Errorf(context.Background(), "Send email fail: %s", err.Error())
		panic(err)
	}
}
