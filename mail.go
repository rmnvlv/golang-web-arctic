package main

import (
	"gopkg.in/gomail.v2"
)

const (
	EmailSubject              = "Thank you for registering for AMTC 2022"
	EmailSubjectMailing       = "AMTC 2022"
	EmailSubjectUploaded      = "AMTC 2022 File uploaded!"
	EmailRegistrationTemplate = `
		<html>
		<body>
		<h3>%s, Thank you for registering at the International Conference «Arctic: Marine Transportation Challenges – 2022» on November 24-25, 2022. </h3>
		<h3>You can find up-to-date information about the key dates of the Conference <a href="%s/programme-overview">here</a>.</h3>
		<h3>If you have any questions, please contact by <a href="mailto:amtc@gumrf.ru">amtc@gumrf.ru</a>.</h3>
		<a href="%s"></a>
		</body>
		</html>
	`
	EmailAbstractsTemplate = `
	<html>
	<body>
	<h3>Dear %s, </h3>
	<h3>You received this email because you are 
	registered for the International Conference «Arctic: Marine Transportation Challenges – 2022», 
	which will take place on November 24-25, 2022.</h3>
	<h2>Acceptance of abstracts is open. The abstracts template is available on the website page</h2>
	<h2>Abstract upload form at the <a href="%s"> link</a>.</h2>
	<h3>If you have any questions, please contact by <a href="mailto:amtc@gumrf.ru">amtc@gumrf.ru</a>.</h3>
	<h4>Best Regards,
	Organizing committee AMTC-2022</h4>
	</body>
	</html>`

	EmailArticleTemplate = `
	<html>
	<body>
	<h3>Dear %s, </h3>
	<h3>You received this email because you are registered for the International Conference 
	«Arctic: Marine Transportation Challenges – 2022»,
	 which will take place on November 24-25, 2022.</h3>
	<h3>Acceptance of full paper is open. The full paper template is available on the website page.</h3>
	<h2>The form for adding full paper is available at the <a href="%s"> link</a>.</h2>
	<h3>If you have any questions, please contact by <a href="mailto:amtc@gumrf.ru">amtc@gumrf.ru</a>.</h3>
	<h3>Best Regards,
	Organizing committee AMTC-2022</h3>
	</body>
	</html>`

	EmailMailingArticleTemplate = `
	<html>
	<body>
	<h3>%s, Full paper uploaded successfully. We will contact you if there are questions about the results of the review. </h3>
	<h3>Please clarify by amtc@gumrf.ru whether an oral presentation is planned or only publication. 
	In the case of an oral presentation, whether it will be a face-to-face or online participation.</h3>
	<h3>If you have any questions, please contact by <a href="mailto:amtc@gumrf.ru">amtc@gumrf.ru</a>.</h3>
	</body>
	</html>
	`
	EmailMailingAbstractsTemplate = `
	<html>
	<body>
	<h3>%s, Abstracts uploaded successfully. </h3>
	<h3>The International Conference «Arctic: Marine Transportation Challenges – 2022» will be held on November 24-25, 2022</h3>
	<h3>If you have any questions, please contact by <a href="mailto:amtc@gumrf.ru">amtc@gumrf.ru</a>.</h3>
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
	//TODO: add email fron .env
	email.SetAddressHeader("From", Cfg["SMTP_USER"], "AMTC 2022")
	email.SetAddressHeader("To", to.Email, to.Name)
	email.SetHeader("Subject", message.Subject)
	email.SetBody("text/html", message.Text)

	return a.mailer.Send(Cfg["SMTP_USER"] /*bad things can happen in here*/, []string{to.Email}, email)
}

// func (a *App) sendNewsletter(mailList []To, message Message) {
// 	for _, to := range mailList {
// 		email := gomail.NewMessage(gomail.SetCharset("UTF-8"), gomail.SetEncoding(gomail.Base64))
// 		email.SetAddressHeader("From", "amtc@gumrf.ru", "AMTC 2022")
// 		email.SetAddressHeader("To", to.Email, to.Name)
// 		email.SetHeader("Subject", message.Subject)
// 		email.SetBody("text/html", message.Text)

// 		if err := a.mailer.Send(os.Getenv("SMTP_USER"), []string{to.Email}, email); err != nil {
// 			a.log.Errorf("Could not send email to %q: %v", to.Email, err)
// 		}

// 		email.Reset()
// 	}
// }
