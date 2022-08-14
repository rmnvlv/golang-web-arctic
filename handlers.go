package main

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net/mail"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize"
	emailverifier "github.com/AfterShip/email-verifier"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

const (
	ErrorMessage   = "Some form fields are entered incorrectly. Please change them."
	SuccessMessage = "Thank you for registration for AMTC 2022!"
)

// var mailQueue []Participant = nil

func (a *App) registerNewParticipant(c *fiber.Ctx) error {
	participant := Participant{
		Surname:             c.FormValue("surname"),
		Name:                c.FormValue("name"),
		Organization:        c.FormValue("organization"),
		Position:            c.FormValue("position"),
		Phone:               c.FormValue("phone"),
		Email:               c.FormValue("email"),
		PresentationForm:    c.FormValue("presentation-form"),
		PresentationSection: c.FormValue("presentation-section"),
		PresentationTitle:   c.FormValue("presentation-title"),
	}

	formErrors := make(map[string]string)

	if os.Getenv("APP_ENV") == "prod" {
		hCaptcha := c.FormValue("h-captcha-response")

		if ok, err := verifyCaptcha(hCaptcha); !ok {
			if errors.Is(err, ErrCaptchaEmpty) {
				formErrors["Captcha"] = "Сaptcha is not passed"
				// formError.Captcha = "Сaptcha is not passed"
			} else {
				// formError.Captcha = "Please try again"
				formErrors["Captcha"] = "Please try again"
			}
			log.Println(err)
		}
	}

	// Validate phone
	val, err := regexp.MatchString(`^((8|\+7)[\- ]?)?(\(?\d{3}\)?[\- ]?)?[\d\- ]{7,10}$`, participant.Phone)
	if err != nil && participant.Phone != "" || !val && participant.Phone != "" {
		formErrors["Phone"] = "Phone number should be valid format."
	}
	//Validate surname
	val, err = regexp.MatchString(`^[a-zA-Z]+$`, participant.Surname)
	if err != nil || !val {
		formErrors["Surname"] = "Surname can only be a-zA-Z."
	}
	//validate name
	val, err = regexp.MatchString(`^[a-zA-Z]+$`, participant.Name)
	if err != nil || !val {
		formErrors["Name"] = "Name can only be a-zA-Z."
	}
	//validate email
	if _, err := mail.ParseAddress(participant.Email); err != nil {
		formErrors["Email"] = "Wrong email format. Example: mail@example.com"
	}
	// val, err = regexp.MatchString(`[^@\s]+@[^@\s]+\.[^@\s]+$`, participant.Email)
	// if err != nil || !val {
	// 	formErrors["Email"] = "Wrong email format. Example: mail@example.com"
	// }
	verifier := emailverifier.NewVerifier()
	if _, err = verifier.Verify(participant.Email); err != nil {
		formErrors["Email"] = "Email does not exists."
	}

	data := fiber.Map{}
	messages := make(map[string]string)

	participant.Code = uuid.New().String()

	if len(formErrors) > 0 {
		messages["Error"] = ErrorMessage
		data["Values"] = participant
	} else {
		a.db.Create(&participant)

		if err := a.sendEmail(
			To{strings.Join([]string{participant.Name, participant.Surname}, " "), participant.Email},
			Message{EmailSubject, EmailAbstractsTemplate},
		); err != nil {
			// log
			fmt.Println(err)
		}

		// mailQueue = append(mailQueue, participant)

		// ch := make(chan []Participant, 1)
		// go worker(ch)
		// ch <- mailQueue
		// close(ch)

		messages["Success"] = SuccessMessage

	}

	data["Title"] = "Registration and submission"
	data["Links"] = Links
	data["Errors"] = formErrors
	data["Message"] = messages

	return c.Render("registration", data)
}

func (a *App) createExcelFile() (*bytes.Buffer, error) {
	var participants []Participant

	result := a.db.Find(&participants)
	if err := result.Error; err != nil {
		return &bytes.Buffer{}, err
	}

	headers := []string{
		"Id",
		"Name",
		"Surname",
		"Organization",
		"Position",
		"Phone",
		"Email",
		"Presentation Form",
		"Presentation Section",
		"Presentation Title",
		"Code",
	}

	document := excelize.NewFile()

	sheetName := "AMTC_2022_Participants"
	_ = document.NewSheet(sheetName)
	document.DeleteSheet("Sheet1")

	for i, h := range headers {
		cellId := string([]byte{uint8(65 + i)}) + "1"
		fmt.Println(cellId)
		document.SetCellValue(sheetName, cellId, h)
	}

	for i, participan := range participants {
		rowIndex := strconv.Itoa(i + 2)
		document.SetCellValue(sheetName, "A"+rowIndex, participan.Name)
		document.SetCellValue(sheetName, "B"+rowIndex, participan.Surname)
		document.SetCellValue(sheetName, "C"+rowIndex, participan.Organization)
		document.SetCellValue(sheetName, "D"+rowIndex, participan.Position)
		document.SetCellValue(sheetName, "E"+rowIndex, participan.Phone)
		document.SetCellValue(sheetName, "F"+rowIndex, participan.Email)
		document.SetCellValue(sheetName, "H"+rowIndex, participan.PresentationSection)
		document.SetCellValue(sheetName, "I"+rowIndex, participan.PresentationTitle)
		document.SetCellValue(sheetName, "J"+rowIndex, participan.Code)
	}

	var buf bytes.Buffer
	if err := document.Write(&buf); err != nil {
		return &bytes.Buffer{}, err
	}

	return &buf, nil
}

func (a *App) downloadFiles(c *fiber.Ctx) error {
	fileType := c.Params("file")

	a.log.Debug(fileType)

	var (
		file     *bytes.Buffer
		err      error
		fileName string
	)

	switch fileType {
	case "participants":
		file, err = a.createExcelFile()
		fileName = fmt.Sprintf("AMTC_2022_%s.%s", "Paticipants", "xlsx")
	case "articles":
		file, err = createZipArchive(a.config.DiskPath + "/" + fileType)
		fileName = fmt.Sprintf("AMTC_2022_%s.%s", "Articles", "zip")
	case "tezisi":
		file, err = createZipArchive(a.config.DiskPath + "/" + fileType)
		fileName = fmt.Sprintf("AMTC_2022_%s.%s", "Tezisi", "zip")
	// case "all":
	// 	file, err = createZipArchive(a.config.DiskPath)
	// 	fileName = fmt.Sprintf("AMTC_2022_%s.%s", "Tesisi+Articles+Participants", "zip")
	default:
		return c.RedirectToRoute("/admin", fiber.Map{
			"Errors": []string{
				fmt.Sprintf("File type doesn't exist: %s", fileType),
			},
		})
	}

	if err != nil {
		return c.RedirectToRoute("admin", fiber.Map{
			"Errors": []string{
				fmt.Sprintf("Can't create zip archive: %s", err.Error()),
			},
		})
	}
	c.Set("Content-Description", "File Transfer")
	c.Set("Content-Disposition", "attachment; filename="+fileName)
	c.Status(fiber.StatusOK)
	return c.SendStream(file)
}

func createZipArchive(src string) (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)

	zw := zip.NewWriter(buf)
	defer zw.Close()

	walker := func(path string, info os.FileInfo, err error) error {
		fmt.Printf("Crawling: %#v\n", path)
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		// Ensure that `path` is not absolute; it should not start with "/".
		// This snippet happens to work because I don't use
		// absolute paths, but ensure your real-world code
		// transforms path into a zip-root relative path.
		p := strings.TrimLeft(path, src)
		fmt.Println(path, p)

		f, err := zw.Create(p)
		if err != nil {
			return err
		}

		_, err = io.Copy(f, file)
		if err != nil {
			return err
		}

		return nil
	}
	if err := filepath.Walk(src, walker); err != nil {
		return nil, fmt.Errorf("walk: %w", err)
	}

	return buf, nil
}

func (a *App) mainView(c *fiber.Ctx) error {
	data := IndexPage
	data["Links"] = Links
	data["Header"] = true
	return c.Render("index", data)
}

func (a *App) programOverviewView(c *fiber.Ctx) error {
	data := fiber.Map{}
	data["Links"] = Links
	data["Title"] = "Programme Overview"
	return c.Render("programm-overview", data)
}

func (a *App) keynoteSpeakersView(c *fiber.Ctx) error {
	data := fiber.Map{}
	data["Title"] = "Keynote Speakers"
	data["Links"] = Links
	data["Content"] = "Key speakers to be determined later."
	return c.Render("basic", data)
}

func (a *App) requirementsView(c *fiber.Ctx) error {
	data := fiber.Map{}
	data["Title"] = "Requirements"
	data["Links"] = Links
	data["Content"] = "Article template will be posted later."
	return c.Render("basic", data)
}

func (a *App) generalInfoView(c *fiber.Ctx) error {
	data := fiber.Map{}
	data["Links"] = Links
	return c.Render("general-information", data)
}

func (a *App) registrationView(c *fiber.Ctx) error {
	data := fiber.Map{}
	data["Title"] = "Registration and submission"
	data["Links"] = Links
	return c.Render("registration", data)
}

func (a *App) adminView(c *fiber.Ctx) error {
	var participants []Participant
	a.db.Find(&participants)
	data := fiber.Map{}
	data["Users"] = participants
	data["Title"] = "Admin"
	data["Links"] = Links
	data["Content"] = "Admin panel"
	return c.Render("admin", data)
}

var form = map[string]map[string]string{
	"tezis": map[string]string{
		"label": "Upload Abstracts",
		"id":    "tezis",
	},
	"article": map[string]string{
		"label": "Upload Full paper",
		"id":    "article",
	},
}

func (a *App) uploadView(c *fiber.Ctx) error {
	// t is type of file: article/tezis
	t := c.Params("type")
	if t != "article" && t != "tezis" {
		a.log.Debug("file type", t)
		return c.Redirect("/404")
	}

	id := c.Query("code")
	if id == "" {
		a.log.Debug("User id is empty")
		return c.Redirect("/404")
	}

	a.log.Debug("User id", id)

	var person Participant
	result := a.db.First(&person, "code = ?", id)
	if result.Error != nil {
		a.log.Error(result.Error)
		return c.Redirect("/404")
	}

	data := fiber.Map{}
	data["Title"] = "Upload"
	data["User"] = person
	data["Form"] = form[t]
	data["Path"] = t + "?code=" + person.Code

	return c.Render("upload", data)
}

const UploadErrorMessage = "Can't upload file."

func (a *App) uploadFile(c *fiber.Ctx) error {
	// t is type of file: article/tezis
	t := c.Params("type")
	var emailMessage string
	if t != "article" && t != "tezis" {
		a.log.Info("file type: ", t)
		return c.Redirect("/404")
	} else if t == "article" {
		emailMessage = EmailMailingArticleTemplate
	} else if t == "tezis" {
		emailMessage = EmailMailingAbstractsTemplate
	}

	a.log.Debug("file type: ", t)

	code := c.Query("code")
	if code == "" {
		a.log.Info("User code is empty")
		return c.Redirect("/404")
	}

	a.log.Debug("User id: ", code)

	var person Participant
	result := a.db.First(&person, "code = ?", code)
	if result.Error != nil {
		a.log.Error(result.Error)
		return c.Redirect("/404")
	}

	data := fiber.Map{}
	data["Title"] = "Upload"
	data["Form"] = form[t]
	data["User"] = person
	data["Path"] = t + "?code=" + person.Code

	file, err := c.FormFile(t)
	if err != nil {
		a.log.Error(err)
		data["Error"] = UploadErrorMessage
		return c.Render("upload", data)
	}

	s := strings.Split(file.Filename, ".")
	ext := s[len(s)-1]

	a.log.Info("file extetton", ext)

	content, err := file.Open()
	if err != nil {
		a.log.Error(err)
		data["Error"] = UploadErrorMessage
		return c.Render("upload", data)
	}

	defer content.Close()

	err = a.saveToDisk(content, ext)
	if err != nil {
		a.log.Errorf("Can't save file to disk: %v", err)

		data["Error"] = UploadErrorMessage
		return c.Render("upload", data)
	}

	data["Success"] = "File successfully uploaded"

	if err = a.sendEmail(
		To{strings.Join([]string{person.Name, person.Surname}, " "), person.Email},
		Message{EmailSubjectUploaded, fmt.Sprintf(emailMessage, strings.Join([]string{person.Name, person.Surname}, " "))},
	); err != nil {
		a.log.Error(err)
	}

	return c.Render("upload", data)
}

func (a *App) sendMailing(c *fiber.Ctx) error {
	var participants []Participant

	a.db.Find(&participants)

	fileForm := c.FormValue("file-form")

	errorEmails := make([]string, 1)
	flag := false

	for _, participant := range participants {
		hrefUpload := fmt.Sprintf("http://%s/%s/%s?code=%s", a.config.Domain, "upload", fileForm, participant.Code)

		nameSurname := strings.Join([]string{participant.Name, participant.Surname}, " ")

		var template string
		switch fileForm {
		case "tezis":
			template = EmailAbstractsTemplate
		case "article":
			template = EmailArticleTemplate
		}

		err := a.sendEmail(
			To{
				nameSurname,
				participant.Email,
			},
			Message{
				EmailSubjectMailing,
				fmt.Sprintf(template, nameSurname, hrefUpload),
			},
		)
		if err != nil {
			a.log.Debug(fmt.Sprintf("Message to email: %s not sent, error: %s", participant.Email, err))
		}

		// time.Sleep(1 * time.Second)

	}

	data := fiber.Map{}

	if flag {
		data["Error"] = fmt.Sprintf("Messages not sent to this emails: %v", errorEmails)
		data["Ending"] = "Message sending completed with some errors"
	} else {
		data["Success"] = "Message sending completed successfully"
		data["Ending"] = "Completed"
	}

	data["Links"] = Links

	return c.Render("admin", data)
}

func (a *App) notFoundView(c *fiber.Ctx) error {
	data := fiber.Map{}
	data["Title"] = "Page Not Found"
	data["Links"] = Links
	data["Content"] = "Page Not Found"
	return c.Render("basic", data)
}
