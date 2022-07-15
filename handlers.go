package main

import (
	"encoding/csv"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

func registerNewParticipant(c *fiber.Ctx) error {
	participant := Participant{
		Surname:             c.FormValue("surname"),
		Name:                c.FormValue("name"),
		Organization:        c.FormValue("organization"),
		Position:            c.FormValue("position"),
		Phone:               c.FormValue("phone"),
		Email:               c.FormValue("email"),
		Type:                c.FormValue("type"),
		Presentation:        c.FormValue("presentation"),
		TitleOfPresentation: c.FormValue("titleofpresentation"),
	}

	var formError Error

	flag := true

	// Validate phone
	val, err := regexp.MatchString(`^((8|\+7)[\- ]?)?(\(?\d{3}\)?[\- ]?)?[\d\- ]{7,10}$`, participant.Phone)
	if err != nil && participant.Phone != "" || !val && participant.Phone != "" {
		formError.Phone = "Invalid phone number"
		flag = false
	}
	//Validate surname
	val, err = regexp.MatchString(`^[a-zA-Z]+$`, participant.Surname)
	if err != nil || !val {
		formError.Surname = "Invalid Surname"
		flag = false
	}
	//validate name
	val, err = regexp.MatchString(`^[a-zA-Z]+$`, participant.Name)
	if err != nil || !val {
		formError.Name = "Invalid Name"
		flag = false
	}
	//validate email
	val, err = regexp.MatchString(`[^@\s]+@[^@\s]+\.[^@\s]+$`, participant.Email)
	if err != nil || !val {
		formError.Email = "Invalid email"
		flag = false
	}

	message := "Registration not completed"

	if flag {
		DB.Create(&participant)
		message = "Registration successefully completed!"
	}

	// после записи в бд форма очищается и в идеале показывается сообщене что операция прошла успешно
	// TODO: показать сообщение об успешной операции или ошибку если такая случилась
	// TODO: create views/registration.html
	return c.Render("registration", fiber.Map{
		"Title":               "Registration and submission",
		"Links":               Content.Links,
		"Error":               formError,
		"Name":                participant.Name,
		"Surname":             participant.Surname,
		"Position":            participant.Position,
		"Organization":        participant.Position,
		"Type":                participant.Type,
		"Presentation":        participant.Presentation,
		"TitleOfPresentation": participant.TitleOfPresentation,
		"Message":             message,
	})
}

func downloadCSV(c *fiber.Ctx) error {
	var participants []Participant

	// Паника почему-то
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
		"Type",
		"Presentation",
		"TitleOfPresentation",
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
			participan.Type,
			participan.Presentation,
			participan.TitleOfPresentation,
		}
		if err := writer.Write(row); err != nil {
			panic(err)
		}
	}

	writer.Flush()

	return c.SendFile(fileName)
}
