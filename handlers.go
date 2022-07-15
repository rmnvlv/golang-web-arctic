package main

import (
	"encoding/csv"
	"os"
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
		Planning:     c.FormValue("planning"),
		Title:        c.FormValue("title"),
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
	// DB.Find(&participants)

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
		"Planning",
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
			participan.Planning,
			participan.Title,
		}
		if err := writer.Write(row); err != nil {
			panic(err)
		}
	}

	writer.Flush()

	return c.SendFile(fileName)
}
