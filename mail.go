package main

import (
	"fmt"

	"gopkg.in/gomail.v2"
)

var (
	auth    AuthInfo
	subject = "Thank you for registration for AMTC 2022"
	message = `
		<html>
		<body>
		<h3>Dear %s, Thank you for taking part in AMTC 2022!</h3>
		<a href="#">Link to site</a>
		</body>
		</html>
	`
)

type To struct {
	Name  string
	Email string
}

type AuthInfo struct {
	Username string `mapstructure:"SMTP_USER"`
	Password string `mapstructure:"SMTP_PASSWORD"`
	Host     string `mapstructure:"SMTP_HOST"`
	Port     int    `mapstructure:"SMTP_PORT"`
}

func SendMail(to To) error {

	email := gomail.NewMessage(gomail.SetCharset("UTF-8"), gomail.SetEncoding(gomail.Base64))
	email.SetAddressHeader("From", auth.Username, "AMTC 2022")
	email.SetAddressHeader("To", to.Email, to.Name)
	email.SetHeader("Subject", subject)
	email.SetBody("text/html", fmt.Sprintf(message, to.Name))

	dialer := gomail.NewDialer(auth.Host, auth.Port, auth.Username, auth.Password)

	return dialer.DialAndSend(email)
}
