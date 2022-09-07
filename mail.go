package main

import (
	"fmt"

	"gopkg.in/gomail.v2"
)

var AfterRegistrationEmail = Message{
	Subject: "Thank you for registering",
	Text: `
		<html>
		<body>
			<p><strong>
				Thank you for registering at the International Conference «Arctic: Marine Transportation Challenges – 2022» on November 24-25, 2022.
			</strong></p>
			<p>You can find up-to-date information about the key dates of the Conference <a href="%s/programme-overview">here</a>.</p>
			<p>If you have any questions, please contact by <a href="mailto:amtc@gumrf.ru">amtc@gumrf.ru</a>.</p>
		</body>
		</html>`,
}
var StartUploadTezisiEmail = Message{
	Subject: "Abstracts upload",
	Text: `
		<html>
		<body>
			<p><strong>Dear %s,</strong></p>
			<p>You received this email because you are registered for the International Conference «Arctic: Marine Transportation Challenges – 2022», which will take place on November 24-25, 2022.</p>
			<p>Acceptance of abstracts is open. The abstracts template is available on the website <a href="%s/programme-overview">page</a>.</p>
			<p>Abstract upload form at the <a href="%s">link</a>.</p>
			<p>If you have any questions, please contact by <a href="mailto:amtc@gumrf.ru">amtc@gumrf.ru</a>.</p>
			<p>Best Regards,<br>Organizing committee AMTC-2022</p>
		</body>
		</html>`,
}
var StartUploadArticlesEmail = Message{
	Subject: "Full paper upload",
	Text: `
		<html>
		<body>
			<p><strong>Dear %s,</strong></p>
			<p>You received this email because you are registered for the International Conference «Arctic: Marine Transportation Challenges – 2022», which will take place on November 24-25, 2022.</p>
			<p>Acceptance of full paper is open. The full paper template is available on the website <a href="%s/programme-overview">page</a>.</p>
			<p>The form for adding full paper is available at the <a href="%s">link</a>.</p>
			<p>If you have any questions, please contact by <a href="mailto:amtc@gumrf.ru">amtc@gumrf.ru</a>.</p>
			<p>Best Regards,<br>Organizing committee AMTC-2022</p>
		</body>
		</html>`,
}
var AfterTezisiUploadEmail = Message{
	Subject: "Abstracts upload",
	Text: `
		<html>
		<body>
			<p><strong>Dear %s, abstracts uploaded successfully.</strong></p>
			<p>The International Conference «Arctic: Marine Transportation Challenges – 2022» will be held on November 24-25, 2022</p>
			<p>If you have any questions, please contact by <a href="mailto:amtc@gumrf.ru">amtc@gumrf.ru</a>.</p>
		</body>
		</html>`,
}
var AfterArticleUploadEmail = Message{
	Subject: "Full paper upload",
	Text: `
		<html>
		<body>
			<p><strong>Dear %s, full paper uploaded successfully.</strong></p>
			<p> We will contact you if there are questions about the results of the review. </p>
			<p>Please clarify by amtc@gumrf.ru whether an oral presentation is planned or only publication. In the case of an oral presentation, whether it will be a face-to-face or online participation.</p>
			<p>If you have any questions, please contact by <a href="mailto:amtc@gumrf.ru">amtc@gumrf.ru</a>.</p>
		</body>
		</html>`,
}
var VerifyEmailMessage = Message{
	Subject: "AMTC 2022 Verify email",
	Text: `
		<html>
		<body>
			<p><strong>Dear participan, please verify your email by this code.</strong></p>
			<p> Code: %s </p>
			<p>If you have any questions, please contact by <a href="mailto:amtc@gumrf.ru">amtc@gumrf.ru</a>.</p>
		</body>
		</html>`,
}

type To struct {
	Name  string
	Email string
}

type Message struct {
	Subject string
	Text    string
}

func (a *App) sendEmail(to To, message Message) error {
	m, err := gomail.NewDialer(a.config.SMTP.Host, a.config.SMTP.Port, a.config.SMTP.User, a.config.SMTP.Password).Dial()
	if err != nil {
		return fmt.Errorf("can't authenticate to an SMTP server: %w", err)
	}
	defer m.Close()

	a.log.Infof("Authenticated to SMTP server: %s:%d", a.config.SMTP.Host, a.config.SMTP.Port)

	email := gomail.NewMessage(gomail.SetCharset("UTF-8"), gomail.SetEncoding(gomail.Base64))
	email.SetAddressHeader("From", a.config.SMTP.User, "AMTC 2022 Organizers")
	email.SetAddressHeader("To", to.Email, to.Name)
	email.SetHeader("Subject", message.Subject)
	email.SetBody("text/html", message.Text)

	return m.Send(a.config.SMTP.User, []string{to.Email}, email)
}
