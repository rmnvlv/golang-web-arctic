package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/gofiber/template/html"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	DB *gorm.DB
)

func main() {
	fmt.Println(time.Now().Unix())
	engine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{
		Views:       engine,
		ViewsLayout: "main",
	})

	DB, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	if err := DB.AutoMigrate(&Participant{}); err != nil {
		panic("failed to migrate database")

	}

	app.Static("/a", "./assets")

	app.Get("/", func(c *fiber.Ctx) error {
		title := "Index"
		return c.Render("index", fiber.Map{
			"Title":   title,
			"Links":   Links,
			"Content": homeContent,
		})
	})

	app.Get("/about", func(c *fiber.Ctx) error {
		title := "About"
		return c.Render("about", fiber.Map{
			"Title":   title,
			"Links":   Links,
			"Content": AboutContent,
		})
	})

	app.Get("/programme-overview", func(c *fiber.Ctx) error {
		title := "Programme Overview"
		content := "The conference program will be posted later."
		return c.Render("basic", fiber.Map{
			"Title":   title,
			"Content": content,
			"Links":   Links,
		})
	})

	app.Get("/keynote-speakers", func(c *fiber.Ctx) error {
		title := "Keynote Speakers"
		content := "Key speakers to be determined later."
		return c.Render("basic", fiber.Map{
			"Title":   title,
			"Content": content,
			"Links":   Links,
		})
	})

	app.Get("/requirements", func(c *fiber.Ctx) error {
		title := "Requirements "
		content := "Article template will be posted later."
		return c.Render("basic", fiber.Map{
			"Title":   title,
			"Links":   Links,
			"Content": content,
		})
	})

	app.Get("/general-information", func(c *fiber.Ctx) error {
		title := "General information "
		return c.Render("general-information", fiber.Map{
			"Title": title,
			"Links": Links,
		})
	})

	app.Get("/registration-and-submission", func(c *fiber.Ctx) error {
		title := "Registration and submission "
		return c.Render("basic", fiber.Map{
			"Title":   title,
			"Links":   Links,
			"Content": "TODO: Registration and submission form",
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
			"Links":   Links,
			"Content": "Admin",
		})
	})
	admin.Get("/csv", downloadCSV)

	app.Use(func(c *fiber.Ctx) error {
		title := "Page Not Found"
		return c.Render("basic", fiber.Map{
			"Title":   title,
			"Content": title,
			"Links":   Links,
		})
	})

	log.Fatal(app.Listen(":8080"))
}
