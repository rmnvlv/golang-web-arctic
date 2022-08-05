package main

import (
	"fmt"
	"os"

	"gopkg.in/gomail.v2"
)

const (
	EmailSubject  = "Thank you for registration for AMTC 2022"
	EmailTemplate = `
		<html>
		<body>
		<h3>%s, Thank you for taking part in AMTC 2022!</h3>
		<a href="%s"></a>
		</body>
		</html>
	`
)

type To struct {
	Name  string
	Email string
}

type Message struct {
	Subject string
	Text    string
}

func (a *App) sendEmail(to To, message Message) error {
	email := gomail.NewMessage(gomail.SetCharset("UTF-8"), gomail.SetEncoding(gomail.Base64))
	email.SetAddressHeader("From", "amtc@gumrf.ru", "AMTC 2022")
	email.SetAddressHeader("To", to.Email, to.Name)
	email.SetHeader("Subject", message.Subject)
	email.SetBody("text/html", fmt.Sprintf(message.Text, to.Name, "link"))

	return a.mailer.Send(os.Getenv("SMTP_USER") /*bad things can happen in here*/, []string{to.Email}, email)
}

func (a *App) sendNewsletter(mailList []To, message Message) error {
	panic("not implemented")
}
