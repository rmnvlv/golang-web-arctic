package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
)

func helloFiber(c *fiber.Ctx) error {

	response := "Hi"

	return c.Send([]byte(response))
}

func main() {
	engine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{
		Views:       engine,
		ViewsLayout: "main",
	})

	app.Static("/a", "./assets")

	app.Get("/", func(c *fiber.Ctx) error {
		title := "Index"
		return c.Render("index", fiber.Map{
			"Title": title,
		})
	})

	app.Get("/about", func(c *fiber.Ctx) error {
		title := "About"
		return c.Render("about", fiber.Map{
			"Title": title,
		})
	})

	app.Get("/programme-overview", func(c *fiber.Ctx) error {
		title := "Programme Overview"
		content := "The conference program will be posted later."
		return c.Render("basic", fiber.Map{
			"Title":   title,
			"Content": content,
		})
	})

	app.Get("/keynote-speakers", func(c *fiber.Ctx) error {
		title := "Keynote Speakers"
		content := "Key speakers to be determined later."
		return c.Render("basic", fiber.Map{
			"Title":   title,
			"Content": content,
		})
	})

	app.Use(func(c *fiber.Ctx) error {
		title := "Page Not Found"
		return c.Render("basic", fiber.Map{
			"Title":   title,
			"Content": title,
		})
	})

	log.Fatal(app.Listen(":8080"))
}
