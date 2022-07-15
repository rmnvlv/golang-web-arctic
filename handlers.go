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
		Surname:      c.FormValue("surname"),
		Name:         c.FormValue("name"),
		Organizacion: c.FormValue("organization"),
		Position:     c.FormValue("position"),
		Phone:        c.FormValue("phone"),
		Email:        c.FormValue("email"),
		Type:         c.FormValue("type"),
		Presentation: c.FormValue("presentation"),
		Title:        c.FormValue("title"),
	}

	// Validate phone
	val, err := regexp.MatchString(`^((8|\+7)[\- ]?)?(\(?\d{3}\)?[\- ]?)?[\d\- ]{7,10}$`, participant.Phone)
	if err != nil && participant.Phone != "" || !val && participant.Phone != "" {
		return c.Send([]byte(err.Error())) // Render error
	}
	//Validate surname
	val, err = regexp.MatchString(`^[a-zA-Z]+$`, participant.Surname)
	if err != nil || !val {
		return c.Send([]byte(err.Error())) // Render error
	}
	//validate name
	val, err = regexp.MatchString(`^[a-zA-Z]+$`, participant.Name)
	if err != nil || !val {
		return c.Send([]byte(err.Error())) // Render error
	}
	//validate email
	val, err = regexp.MatchString(`[^@\s]+@[^@\s]+\.[^@\s]+$`, participant.Email)
	if err != nil || !val {
		return c.Send([]byte(err.Error())) // Render error
	}

	DB.Create(&participant)

	// после записи в бд форма очищается и в идеале показывается сообщене что операция прошла успешно
	// TODO: показать сообщение об успешной операции или ошибку если такая случилась
	// TODO: create views/registration.html
	return c.Redirect("/registration-and-submission")
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
		"Title",
	}

	if err := writer.Write(headers); err != nil {
		panic(err)
	}

	for _, participan := range participants {
		row := []string{
			participan.Name,
			participan.Surname,
			participan.Organizacion,
			participan.Phone,
			participan.Email,
			participan.Type,
			participan.Presentation,
			participan.Title,
		}
		if err := writer.Write(row); err != nil {
			panic(err)
		}
	}

	writer.Flush()

	return c.SendFile(fileName)
}
