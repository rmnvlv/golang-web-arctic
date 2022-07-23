package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html"
	"gorm.io/driver/postgres"
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

	dbURL := "postgres://psg:psg@localhost:5432/psg"

	var err error
	DB, err = gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		fmt.Println("error: ", err)
		panic("failed to connect database")
	}

	if err := DB.AutoMigrate(&Participant{}); err != nil {
		panic("failed to migrate database")

	}
	if err := DB.AutoMigrate(&LoadedFile{}); err != nil {
		panic("failed to migrate database")

	}

	app.Static("/a", "./assets")

	app.Get("/", func(c *fiber.Ctx) error {
		data := IndexPage
		data["Links"] = Links
		data["Header"] = true
		return c.Render("index", data)
	})

	app.Get("/programme-overview", func(c *fiber.Ctx) error {
		data := fiber.Map{}
		data["Links"] = Links
		data["Title"] = "Programme Overview"
		return c.Render("programm-overview", data)
	})

	app.Get("/keynote-speakers", func(c *fiber.Ctx) error {
		data := fiber.Map{}
		data["Title"] = "Keynote Speakers"
		data["Links"] = Links
		data["Content"] = "Key speakers to be determined later."
		return c.Render("basic", data)
	})

	app.Get("/requirements", func(c *fiber.Ctx) error {
		data := fiber.Map{}
		data["Title"] = "Requirements"
		data["Links"] = Links
		data["Content"] = "Article template will be posted later."
		return c.Render("basic", data)
	})

	app.Get("/general-information", func(c *fiber.Ctx) error {
		data := fiber.Map{}
		data["Links"] = Links
		return c.Render("general-information", data)
	})

	app.Get("/registration-and-submission", func(c *fiber.Ctx) error {
		data := fiber.Map{}
		data["Title"] = "Registration and submission"
		data["Links"] = Links
		return c.Render("registration", data)
	})

	app.Post("/registration-and-submission", registerNewParticipant)

	admin := app.Group("/admin", basicauth.New(basicauth.Config{
		Users: map[string]string{
			"admin": os.Getenv("ADMIN_PASSWORD"),
		},
	}))
	admin.Get("/", func(c *fiber.Ctx) error {
		data := fiber.Map{}
		data["Title"] = "Admin"
		data["Links"] = Links
		data["Content"] = "Admin panel"
		return c.Render("admin", data)
	})
	admin.Get("/file", downloadFile)

	app.Use(func(c *fiber.Ctx) error {
		data := fiber.Map{}
		data["Title"] = "Page Not Found"
		data["Links"] = Links
		return c.Render("basic", data)
	})

	log.Fatal(app.Listen(":" + os.Getenv("PORT")))
	// log.Fatal(app.Listen(":8080"))

}
