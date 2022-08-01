package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	DB    *gorm.DB
	DBURL = "test.db" /*os.Getenv("DATABASE_URL")*/
)

func main() {
	var err error
	DB, err = gorm.Open(sqlite.Open(DBURL), &gorm.Config{})
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

	engine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{
		Views:       engine,
		ViewsLayout: "main",
	})

	app.Use(logger.New())
	app.Static("/a", "./assets")
	app.Get("/", mainView)
	app.Get("/programme-overview", programOverviewView)
	app.Get("/keynote-speakers", keynoteSpeakersView)
	app.Get("/requirements", requirementsView)
	app.Get("/general-information", generalInfoView)
	app.Get("/registration-and-submission", registrationView)
	app.Post("/registration-and-submission", registerNewParticipant)
	admin := app.Group("/admin", basicauth.New(basicauth.Config{
		Users: map[string]string{
			"admin": "123456", /*os.Getenv("ADMIN_PASSWORD")*/
		},
	}))
	admin.Get("/", adminView)
	admin.Get("/file", downloadFile)
	app.Get("/upload", uploadView)
	app.Post("/upload", uploadArticleOrTezisi)
	app.Use(notFoundView)

	log.Fatal(app.Listen(":" + "8080" /*os.Getenv("PORT")*/))
}
