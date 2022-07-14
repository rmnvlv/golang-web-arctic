package dynamic

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rmnvlv/Web-Arctic/database"
)

func InitUserRoutes(router fiber.Router) {
	registration := router.Group("/registration")

	registration.Get("", registrationGet)
	registration.Post("", registrationPost)
}

func registrationGet(c *fiber.Ctx) error {
	response := "Implement form of registration"
	c.Get(response)
	return nil
}

func registrationPost(c *fiber.Ctx) error {
	var participant Participant

	if err := c.BodyParser(&participant); err != nil {
		return err
	}

	//Validate?
	//How to stщrage files

	db := database.DBconn

	db.Create(&participant)

	c.JSON(participant)

	return nil
}
