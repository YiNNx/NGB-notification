package util

import (
	"github.com/jordan-wright/email"
	"net/smtp"
	"net/textproto"
	"ngb-noti/config"
)

var (
	host     = config.C.Mail.Host
	addr     = config.C.Mail.Addr
	username = config.C.Mail.Username
	password = config.C.Mail.Password
)

func SendEmail(to string, subject string, text string) error {

	e := &email.Email{
		To:      []string{to},
		From:    username,
		Subject: subject,
		Text:    []byte(text),
		Headers: textproto.MIMEHeader{},
	}

	err := e.Send(addr, smtp.PlainAuth("", username, password, host))
	if err != nil {
		return err
	}
	return nil
}
