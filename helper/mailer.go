package helper

import (
	"greenenvironment/configs"
	"net/smtp"
)

type MailerInterface interface {
	SendEmail(smtpConfig configs.SMTPConfig, to string, subject string, body string) error
}

func SendEmail(smtpConfig configs.SMTPConfig, to string, subject string, body string) error {
	auth := smtp.PlainAuth("", smtpConfig.Username, smtpConfig.Password, smtpConfig.Host)

	msg := []byte("Subject: " + subject + "\r\n\r\n" + body)

	err := smtp.SendMail(
		smtpConfig.Host+":"+smtpConfig.Port,
		auth,
		smtpConfig.Username,
		[]string{to},
		msg,
	)

	return err
}