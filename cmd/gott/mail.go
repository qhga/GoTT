package main

import (
	"fmt"
	"net/smtp"
)

func sendMail(to string, subject string, mailTxt string) error {
	from := config.MailUser
	header := `From: TT <%s>
To: %s
Subject: %s


`
	header = fmt.Sprintf(header, from, to, subject)
	mailTxt = header + mailTxt
	auth := smtp.PlainAuth("", from, config.MailPass, config.MailSmtpServer)
	err := smtp.SendMail(config.MailSmtpServer+":"+config.MailSmtpPort, auth, from, []string{to}, []byte(mailTxt))
	return err
}
