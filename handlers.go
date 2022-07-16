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

	if ok, _ := verifyCaptcha(hCaptcha); !ok {
		formError.Captcha = "Please try again"
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

	return c.SendFile(fileName)
}
