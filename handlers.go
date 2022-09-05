package main

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/mail"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	emailverifier "github.com/AfterShip/email-verifier"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

const (
	ErrorMessage   = "Some form fields are entered incorrectly. Change them and try again."
	SuccessMessage = "Seccessfully registered for AMTC 2022!"
)

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

	participant.CreatedAt = time.Now().Format("01-02-2002")

	formErrors := make(map[string]string)

	if a.config.Captcha.Enable {
		hCaptcha := c.FormValue("h-captcha-response")

		if ok, err := verifyCaptcha(hCaptcha); !ok {
			if errors.Is(err, ErrCaptchaEmpty) {
				formErrors["Captcha"] = "Сaptcha is not passed"
			} else {
				formErrors["Captcha"] = "Please try again"
			}
			a.log.Error(err.Error())
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

	verifier := emailverifier.NewVerifier()
	if _, err = verifier.Verify(participant.Email); err != nil {
		formErrors["Email"] = "Email does not exists."
	}

	data := fiber.Map{}
	messages := make(map[string]string)

	participant.Token = uuid.New().String()

	if len(formErrors) > 0 {
		messages["Error"] = ErrorMessage
		data["Values"] = participant
	} else {
		a.db.Create(&participant)

		a.log.Debug(a.db.First(&participant, participant.Token))

		if err := a.sendEmail(
			To{strings.Join([]string{participant.Name, participant.Surname}, " "), participant.Email},
			Message{AfterRegistrationEmail.Subject, fmt.Sprintf(AfterRegistrationEmail.Text, a.config.Domain)},
		); err != nil {
			a.log.Errorf("Can't send email to %s: %w", participant.Email, err.Error())
		}

		messages["Success"] = SuccessMessage
	}

	data["Title"] = "Registration and submission"
	data["Errors"] = formErrors
	data["Message"] = messages

	return c.Render("registration", data)
}

func (a *App) createExcelFile() (*bytes.Buffer, error) {
	var participants []Participant

	result := a.db.Find(&participants)
	a.log.Debug(result)
	if err := result.Error; err != nil {
		return &bytes.Buffer{}, err
	}

	headers := []string{
		"Registred at",
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
		document.SetCellValue(sheetName, "A"+rowIndex, participan.CreatedAt)
		document.SetCellValue(sheetName, "B"+rowIndex, participan.Name)
		document.SetCellValue(sheetName, "C"+rowIndex, participan.Surname)
		document.SetCellValue(sheetName, "D"+rowIndex, participan.Organization)
		document.SetCellValue(sheetName, "E"+rowIndex, participan.Position)
		document.SetCellValue(sheetName, "F"+rowIndex, participan.Phone)
		document.SetCellValue(sheetName, "G"+rowIndex, participan.Email)
		document.SetCellValue(sheetName, "H"+rowIndex, participan.PresentationForm)
		document.SetCellValue(sheetName, "I"+rowIndex, participan.PresentationSection)
		document.SetCellValue(sheetName, "J"+rowIndex, participan.PresentationTitle)
		document.SetCellValue(sheetName, "K"+rowIndex, participan.Token)
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
	case "article":
		file, err = createZipArchive(a.config.DiskPath + "/" + fileType)
		fileName = fmt.Sprintf("AMTC_2022_%s.%s", "Articles", "zip")
	case "tezis":
		file, err = createZipArchive(a.config.DiskPath + "/" + fileType)
		fileName = fmt.Sprintf("AMTC_2022_%s.%s", "Tezisi", "zip")
	case "open-upload":
		file, err = createZipArchive(a.config.DiskPath + "/" + fileType)
		fileName = fmt.Sprintf("AMTC_2022_%s.%s", "Open-upload", "zip")
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
	//result := a.db.First(&person, "token = ?", id)
	result := a.db.Where("token = ?", id).First(&person)
	if result.Error != nil {
		a.log.Error(result.Error)
		return c.Redirect("/404")
	}

	data := fiber.Map{}
	data["Title"] = "Upload"
	data["User"] = person
	data["Form"] = form[t]
	data["Path"] = t + "?code=" + person.Token

	return c.Render("upload", data)
}

const UploadErrorMessage = "Can't upload file."

func (a *App) uploadFile(c *fiber.Ctx) error {
	// t is type of file: article/tezis
	t := c.Params("type")
	a.log.Debug("file type: ", t)

	var emailMessage Message
	if t != "article" && t != "tezis" {
		a.log.Info("file type: ", t)
		return c.Redirect("/404")
	} else if t == "article" {
		emailMessage = AfterArticleUploadEmail
	} else if t == "tezis" {
		emailMessage = AfterTezisiUploadEmail
	}

	token := c.Query("code")
	if token == "" {
		a.log.Info("User code is empty")
		return c.Redirect("/404")
	}

	a.log.Debug("User id: ", token)

	var person Participant
	//result := a.db.First(&person, "token = ?", token)
	result := a.db.Where("token = ?", token).First(&person)
	if result.Error != nil {
		a.log.Error(result.Error)
		return c.Redirect("/404")
	}

	data := fiber.Map{}
	data["Title"] = "Upload"
	data["Form"] = form[t]
	data["User"] = person
	data["Path"] = t + "?code=" + person.Token

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

	fileName := fmt.Sprintf("%s_%s_%s_%s", person.Name, person.Surname, person.Email, t)

	err = a.saveToDisk(content, ext, t+"/"+fileName)
	if err != nil {
		a.log.Errorf("Can't save file to disk: %v", err)

		data["Error"] = UploadErrorMessage
		return c.Render("upload", data)
	}

	data["Success"] = "File successfully uploaded"

	if err = a.sendEmail(
		To{strings.Join([]string{person.Name, person.Surname}, " "), person.Email},
		Message{emailMessage.Subject, fmt.Sprintf(emailMessage.Text, strings.Join([]string{person.Name, person.Surname}, " "))},
	); err != nil {
		a.log.Error(err)
	}

	return c.Render("upload", data)
}

func (a *App) sendNewsletter(c *fiber.Ctx) error {
	var participants []Participant

	if err := a.db.Find(&participants).Error; err != nil {
		return c.RedirectToRoute("/admin", fiber.Map{"Errors": map[string]string{"sendNewsletter": "Can't get participants"}})
	}

	fileForm := c.FormValue("file-form")

	errorEmails := make([]string, 1)
	flag := false

	for _, participant := range participants {
		hrefUpload := fmt.Sprintf("%s/upload/%s?code=%s", a.config.Domain, fileForm, participant.Token)

		nameSurname := strings.Join([]string{participant.Name, participant.Surname}, " ")

		var emailTemplate Message
		switch fileForm {
		case "tezis":
			emailTemplate = StartUploadTezisiEmail
		case "article":
			emailTemplate = StartUploadTezisiEmail
		}

		err := a.sendEmail(
			To{
				nameSurname,
				participant.Email,
			},
			Message{
				emailTemplate.Subject,
				fmt.Sprintf(emailTemplate.Text, nameSurname, hrefUpload),
			},
		)
		if err != nil {
			a.log.Debug(fmt.Sprintf("Message to email: %s not sent, error: %s", participant.Email, err))
			errorEmails = append(errorEmails, participant.Email)
			flag = true
		}
	}

	data := fiber.Map{}

	if flag {
		data["Error"] = fmt.Sprintf("Messages not sent to this emails: %v", errorEmails)
		data["Ending"] = "Message sending completed with some errors"
	} else {
		data["Success"] = "Message sending completed successfully"
		data["Ending"] = "Completed"
	}

	return c.Render("admin", data)
}

func (a *App) openUpload(c *fiber.Ctx) error {
	name := c.FormValue("name")
	surname := c.FormValue("surname")
	email := c.FormValue("email")

	formErrors := make(map[string]string)

	//Validate surname
	val, err := regexp.MatchString(`^[a-zA-Z]+$`, surname)
	if err != nil || !val {
		formErrors["Surname"] = "Surname can only be a-z A-Z."
	}
	//validate name
	val, err = regexp.MatchString(`^[a-zA-Z]+$`, name)
	if err != nil || !val {
		formErrors["Name"] = "Name can only be a-z A-Z."
	}
	//validate email
	if _, err := mail.ParseAddress(email); err != nil {
		formErrors["Email"] = "Wrong email format. Example: mail@example.com"
	}
	verifier := emailverifier.NewVerifier()
	if _, err = verifier.Verify(email); err != nil {
		formErrors["Email"] = "Email does not exists."
	}

	//Upload file

	data := fiber.Map{}
	data["Title"] = "Opened upload"

	messages := make(map[string]string)

	if len(formErrors) > 0 {
		messages["Error"] = ErrorMessage
	} else {

		file, err := c.FormFile("article")
		if err != nil {
			a.log.Error(err)
			messages["Error"] = UploadErrorMessage
			data["Message"] = messages
			return c.Render("open-upload", data)
		}

		s := strings.Split(file.Filename, ".")
		ext := s[len(s)-1]

		a.log.Info("file extetton", ext)

		content, err := file.Open()
		if err != nil {
			a.log.Error(err)
			messages["Error"] = UploadErrorMessage
			data["Message"] = messages
			return c.Render("open-upload", data)
		}
		defer content.Close()

		uniqueId := uuid.New().String()
		fileName := fmt.Sprintf("%s_%s_%s_%s_%s", name, surname, email, "open-upload", uniqueId)

		err = a.saveToDisk(content, ext, "open-upload/"+fileName)
		if err != nil {
			a.log.Errorf("Can't save file to disk: %v", err)
			messages["Error"] = UploadErrorMessage
			data["Message"] = messages
			return c.Render("open-upload", data)
		}

		messages["Success"] = "File successfully uploaded"

		if err := a.sendEmail(
			To{strings.Join([]string{name, surname}, " "), email},
			Message{AfterTezisiUploadEmail.Subject, fmt.Sprintf(AfterArticleUploadEmail.Text, strings.Join([]string{name, surname}, " "))},
		); err != nil {
			a.log.Error(err)
		}
	}

	data["Message"] = messages
	data["Errors"] = formErrors

	return c.Render("open-upload", data)
}

func (a *App) openUploadView(c *fiber.Ctx) error {
	//month-day now
	dtNow := time.Now().Format("01-02")
	monthNow, err := strconv.Atoi(dtNow[0:2])
	if err != nil {
		return err
	}
	dayNow, err := strconv.Atoi(dtNow[3:])
	if err != nil {
		return err
	}
	a.log.Debug(dtNow, monthNow, dayNow)

	//TODO: From .env
	dtNeed := a.config.UploadingDate
	monthNeed, err := strconv.Atoi(dtNeed[0:2])
	if err != nil {
		return err
	}
	dayNeed, err := strconv.Atoi(dtNeed[3:])
	if err != nil {
		return err
	}
	a.log.Debug(dtNeed, monthNeed, dayNeed)

	data := fiber.Map{}

	c.Bind(fiber.Map{
		"Title": "Upload",
	})

	if monthNeed <= monthNow && dayNeed <= dayNow {
		return c.Render("open-upload", data)
	} else {
		//render error
		data["Closed"] = "Uploading articles for unregistered participants is closed for now, come back later"
		a.log.Debug("RENDER CLOSED UPLOAD")

		return c.Render("open-upload", data)
	}
}
