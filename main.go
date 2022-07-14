package main

import (
	"fmt"
	"log"

	"github.com/gofiber/template/html"

	"github.com/gofiber/fiber/v2"
	"github.com/rmnvlv/Web-Arctic/database"
	"github.com/rmnvlv/Web-Arctic/dynamic"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func initDatabase() {
	var err error
	database.DBconn, err = gorm.Open(sqlite.Open("test.db"))

	if err != nil {
		panic("Failed to connect to db")
	}
	fmt.Println("Succsess to connect to db")

	database.DBconn.AutoMigrate(&dynamic.Participant{})
	fmt.Println("DB migratetd")
}

// Init routes to get particions
func initDynamicRoutes(app *fiber.App) {

	admin := app.Group("/admin")
	dynamic.InitAdminRouter(admin)

	user := app.Group("/user")
	dynamic.InitUserRoutes(user)
}

func main() {
	engine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{
		Views:       engine,
		ViewsLayout: "main",
	})

	initDatabase()
	initDynamicRoutes(app)

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
