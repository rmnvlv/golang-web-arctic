package dynamic

import (
	"encoding/csv"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/rmnvlv/Web-Arctic/database"
)

func InitAdminRouter(router fiber.Router) {
	particions := router.Group("/participants")
	particions.Get("", downloadParticipantsGet)
	particions.Post("", downloadParticipantsPost)
}

func downloadParticipantsGet(c *fiber.Ctx) error {
	response := "Implement form of download and button"

	c.Send([]byte(response))

	return nil
}

func downloadParticipantsPost(c *fiber.Ctx) error {
	var input Administrator

	if err := c.BodyParser(&input); err != nil {
		return err
	}

	if input.Secret != "Fuck_me_for_the_win007" {
		return fiber.NewError(500, "Bad password")
	}

	db := database.DBconn

	var participants []Participant

	db.Find(&participants)

	c.JSON(participants)

	csvFile, err := os.Create("./FileOfParticipants.csv")

	if err != nil {
		return err
	}

	defer csvFile.Close()

	writer := csv.NewWriter(csvFile)

	for _, participan := range participants {
		var row []string
		row = append(row, participan.Name)
		row = append(row, participan.Surname)
		row = append(row, participan.Organizacion)
		row = append(row, participan.Phone)
		row = append(row, participan.Email)
		row = append(row, participan.Type)
		row = append(row, participan.Title)
		writer.Write(row)
	}

	writer.Flush()

	c.SendFile("./FileOfParticipants.csv")

	return nil
}
