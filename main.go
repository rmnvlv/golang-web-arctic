package main

import "github.com/gofiber/fiber/v2"

func helloFiber(c *fiber.Ctx) error {

	response := "Hi"

	c.Send([]byte(response))

	return nil
}

func main() {
	app := fiber.New()

	app.Get("/", helloFiber)

	app.Listen(":8080")
}
