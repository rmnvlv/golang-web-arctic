package main

import (
	"fmt"
	"os"

	"gopkg.in/gomail.v2"
)

const (
	EmailSubject              = "Thank you for registering for AMTC 2022"
	EmailRegistrationTemplate = `
		<html>
		<body>
		<h3>%s, Thank you for registering at the International Conference «Arctic: Marine Transportation Challenges – 2022» on November 24-25, 2022. </h3>
		<h3>You can find up-to-date information about the key dates of the Conference here.</h3>
		<h3>If you have any questions, please contact by amtc@gumrf.ru.</h3>
		<a href="%s"></a>
		</body>
		</html>
	`
	EmailAbstractsArticleTemplate = `
	<<html>
	<body>
	<h3>Dear %s, </h3>
	<h3>You received this email because you are registered for the International Conference 
	«Arctic: Marine Transportation Challenges – 2022», which will take place on November 24-25, 2022.
	Acceptance of full paper and abstracts is open. The full paper and abstracts template is available on the website page.</h3>
	<h2>Abstract upload form at the link.</h2>
	<h2>The form for adding full paper is available at the link.</h2><a href="%s"></a>
	<h3>If you have any questions, please contact by amtc@gumrf.ru.</h3><a href="%s"></a>
	<a href="%s"></a>
	<h4>Best Regards,</h4>
	<h4>Organizing committee AMTC-2022</h4>
	</body>
	</html>>`

	EmailMailingArticleTemplate = `
	<html>
	<body>
	<h3>%s, Full paper uploaded successfully. We will contact you if there are questions about the results of the review. </h3>
	<h3>Please clarify by amtc@gumrf.ru whether an oral presentation is planned or only publication. 
	In the case of an oral presentation, whether it will be a face-to-face or online participation.</h3>
	<h3>If you have any questions, please contact by amtc@gumrf.ru.</h3>
	<a href="%s"></a>
	</body>
	</html>
	`
	EmailMailingAbstractsTemplate = `
	<html>
	<body>
	<h3>%s, Abstracts uploaded successfully. </h3>
	<h3>The International Conference «Arctic: Marine Transportation Challenges – 2022» will be held on November 24-25, 2022</h3>
	<h3>If you have any questions, please contact by amtc@gumrf.ru.</h3>
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
