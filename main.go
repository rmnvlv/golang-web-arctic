package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	DB *gorm.DB
)

func main() {
	engine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{
		Views:       engine,
		ViewsLayout: "main",
	})

	app.Use(logger.New())

	var err error
	DB, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	if err := DB.AutoMigrate(&Participant{}); err != nil {
		panic("failed to migrate database")

	}

	app.Static("/a", "./assets")

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{
			"Title":   "Conference",
			"Links":   Content.Links,
			"Content": Content.Home,
		})
	})

	app.Get("/about", func(c *fiber.Ctx) error {
		return c.Render("about", fiber.Map{
			"Title":   "About",
			"Links":   Content.Links,
			"Content": Content.About,
		})
	})

	app.Get("/programme-overview", func(c *fiber.Ctx) error {
		title := "Programme Overview"
		content := "The conference program will be posted later."
		return c.Render("basic", fiber.Map{
			"Title":   title,
			"Links":   Content.Links,
			"Content": content,
		})
	})

	app.Get("/keynote-speakers", func(c *fiber.Ctx) error {
		title := "Keynote Speakers"
		content := "Key speakers to be determined later."
		return c.Render("basic", fiber.Map{
			"Title":   title,
			"Links":   Content.Links,
			"Content": content,
		})
	})

	app.Get("/requirements", func(c *fiber.Ctx) error {
		title := "Requirements "
		content := "Article template will be posted later."
		return c.Render("basic", fiber.Map{
			"Title":   title,
			"Links":   Content.Links,
			"Content": content,
		})
	})

	app.Get("/general-information", func(c *fiber.Ctx) error {
		title := "General information "
		return c.Render("general-information", fiber.Map{
			"Title": title,
			"Links": Content.Links,
		})
	})

	app.Get("/registration-and-submission", func(c *fiber.Ctx) error {
		return c.Render("registration", fiber.Map{
			"Title": "Registration and submission",
			"Links": Content.Links,
		})
	})

	app.Post("/registration-and-submission", registerNewParticipant)

	admin := app.Group("/admin", basicauth.New(basicauth.Config{
		Users: map[string]string{
			"admin": "123456", //insecure - TODO: get password from env
		},
	}))
	admin.Get("/", func(c *fiber.Ctx) error {
		title := "Admin"
		return c.Render("admin", fiber.Map{
			"Title":   title,
			"Links":   Content.Links,
			"Content": "Admin",
		})
	})
	admin.Get("/csv", downloadCSV)

	app.Use(func(c *fiber.Ctx) error {
		title := "Page Not Found"
		return c.Render("basic", fiber.Map{
			"Title":   title,
			"Content": title,
			"Links":   Content.Links,
		})
	})

	log.Fatal(app.Listen(":" + os.Getenv("PORT")))
}
