package main

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/gofiber/fiber/v2"
)

const (
	ErrorMessage   = "Some form fields are entered incorrectly. Please change them."
	SuccessMessage = "Thank you for registration!"

	hCaptchaAPIURL = "https://hcaptcha.com/siteverify"
)

var (
	hCaptchaSecretKey = os.Getenv("HCAPTCHA_SECRET_KEY")
	hCaptchaSiteKey   = os.Getenv("HCAPTCHA_SITE_KEY")
)

// Response is the hcaptcha JSON response.
type Response struct {
	ChallengeTS string   `json:"challenge_ts"`
	Hostname    string   `json:"hostname"`
	ErrorCodes  []string `json:"error-codes,omitempty"`
	Success     bool     `json:"success"`
	Credit      bool     `json:"credit,omitempty"`
}

func verifyCaptcha(captcha string) (bool, error) {
	if captcha == "" {
		return false, errors.New("captcha is empty")
	}

	form := url.Values{}
	form.Add("secret", hCaptchaSecretKey)
	form.Add("response", captcha)
	form.Add("sitekey", hCaptchaSiteKey)

	resp, err := http.DefaultClient.PostForm(hCaptchaAPIURL, form)
	if err != nil {
		return false, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	defer resp.Body.Close()

	var response Response
	if err := json.Unmarshal(body, &response); err != nil {
		return false, err
	}

	if !response.Success {
		return false, fmt.Errorf("hCaptcha: %v", response.ErrorCodes)
	}

	return true, nil
}

func registerNewParticipant(c *fiber.Ctx) error {
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

	hCaptcha := c.FormValue("h-captcha-response")

	var formError FormError

	if ok, err := verifyCaptcha(hCaptcha); !ok {
		formError.Captcha = "Please try again"
		fmt.Println(err)
	}

	// Validate phone
	val, err := regexp.MatchString(`^((8|\+7)[\- ]?)?(\(?\d{3}\)?[\- ]?)?[\d\- ]{7,10}$`, participant.Phone)
	if err != nil && participant.Phone != "" || !val && participant.Phone != "" {
		formError.Phone = "Phone number should be valid format."
	}
	//Validate surname
	val, err = regexp.MatchString(`^[a-zA-Z]+$`, participant.Surname)
	if err != nil || !val {
		formError.Surname = "Surname can only be a-zA-Z."
	}
	//validate name
	val, err = regexp.MatchString(`^[a-zA-Z]+$`, participant.Name)
	if err != nil || !val {
		formError.Name = "Name can only be a-zA-Z."
	}
	//validate email
	val, err = regexp.MatchString(`[^@\s]+@[^@\s]+\.[^@\s]+$`, participant.Email)
	if err != nil || !val {
		formError.Email = "Wrong email format. Example: maria@example.com."
	}

	var message string
	var formData Participant = participant

	if reflect.DeepEqual(formError, FormError{}) {
		DB.Create(&participant)
		message = SuccessMessage
		formData = Participant{}
	} else {
		formError.Message = ErrorMessage
	}

	return c.Render("registration", fiber.Map{
		"Title":    "Registration and submission",
		"Links":    Content.Links,
		"Error":    formError,
		"FormData": formData,
		"Message":  message,
	})
}

func downloadCSV(c *fiber.Ctx) error {
	var participants []Participant

	DB.Find(&participants)

	fileName := "./" + strconv.FormatInt(time.Now().Unix(), 10) + ".csv"
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)

	headers := []string{
		"Name",
		"Surname",
		"Organization",
		"Position",
		"Phone",
		"Email",
		"Presentation Form",
		"Presentation Section",
		"Presentation Title",
	}

	if err := writer.Write(headers); err != nil {
		panic(err)
	}

	for _, participan := range participants {
		row := []string{
			participan.Name,
			participan.Surname,
			participan.Organization,
			participan.Position,
			participan.Phone,
			participan.Email,
			participan.PresentationForm,
			participan.PresentationSection,
			participan.PresentationTitle,
		}
		if err := writer.Write(row); err != nil {
			panic(err)
		}
	}

	writer.Flush()

	fileExcel := excelize.NewFile()

	_ = fileExcel.NewSheet("Sheet1")

	fileExcel.SetCellValue("Sheet1", "A1", "Name")
	fileExcel.SetCellValue("Sheet1", "B1", "Surname")
	fileExcel.SetCellValue("Sheet1", "C1", "Organization")
	fileExcel.SetCellValue("Sheet1", "D1", "Position")
	fileExcel.SetCellValue("Sheet1", "E1", "Phone")
	fileExcel.SetCellValue("Sheet1", "F1", "Email")
	fileExcel.SetCellValue("Sheet1", "G1", "PresentationForm")
	fileExcel.SetCellValue("Sheet1", "H1", "PresentationSection")
	fileExcel.SetCellValue("Sheet1", "I1", "PresentationTitle")

	// "ABCDEFGHI"
	for counter, participan := range participants {
		number := strconv.Itoa(counter + 2)
		fileExcel.SetCellValue("Sheet1", "A"+number, participan.Name)
		fileExcel.SetCellValue("Sheet1", "B"+number, participan.Surname)
		fileExcel.SetCellValue("Sheet1", "C"+number, participan.Organization)
		fileExcel.SetCellValue("Sheet1", "D"+number, participan.Position)
		fileExcel.SetCellValue("Sheet1", "E"+number, participan.Phone)
		fileExcel.SetCellValue("Sheet1", "F"+number, participan.Email)
		fileExcel.SetCellValue("Sheet1", "G"+number, participan.PresentationForm)
		fileExcel.SetCellValue("Sheet1", "H"+number, participan.PresentationSection)
		fileExcel.SetCellValue("Sheet1", "I"+number, participan.PresentationTitle)
	}

	fileNameExcel := strconv.FormatInt(time.Now().Unix(), 10) + ".xlsx" //  "./" +

	if err = fileExcel.SaveAs(fileNameExcel); err != nil {
		return err
	}

	return c.SendFile("./" + fileNameExcel)
}

func UploadFile(c *fiber.Ctx) error {
	return nil
}
