package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html"
	"github.com/joho/godotenv"
	"gopkg.in/gomail.v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type App struct {
	server *fiber.App
	db     *gorm.DB
	mailer gomail.SendCloser
	// logger *zap.SugaredLogger
	// disk   *YandexDsik
	// captcha *HCaptcha
}

func NewApp() (*App, error) {
	dbURL := os.Getenv("DATABASE_URL")

	// TODO: .env DB_URL="/name" + pull docker with my sql
	// db, err := gorm.Open(mysql.Open(dbURL), &gorm.Config{})
	// if err != nil {
	// 	return nil, fmt.Errorf("can't open database: %w", err)
	// }

	db, err := gorm.Open(sqlite.Open(dbURL), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("can't open database")
	}

	if err := db.AutoMigrate(&Participant{}); err != nil {
		return nil, fmt.Errorf("can't apply migrations to database: %w", err)
	}

	smtpPort, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		return nil, fmt.Errorf("can't convert an SMTP server port to int: %w", err)
	}
	smtpHost := os.Getenv("SMTP_HOST")
	smtpUser := os.Getenv("SMTP_USER")
	smtpPassword := os.Getenv("SMTP_PASSWORD")

	m, err := gomail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPassword).Dial()
	if err != nil {
		return nil, fmt.Errorf("can't authenticate to an SMTP server: %w", err)
	}

	server := fiber.New(fiber.Config{
		Views:       html.New("./views", ".html"),
		ViewsLayout: "main",
	})

	app := App{db: db, mailer: m, server: server}

	return &app, nil
}

func (a *App) Run() {
	a.server.Listen(":" + os.Getenv("PORT"))
}

func (a *App) Shutdown(_ context.Context) error {
	e := make([]string, 0)

	if err := a.server.Shutdown(); err != nil {
		e = append(e, fmt.Errorf("can't shutdown server: %w", err).Error())
	}

	db, err := a.db.DB()
	if err != nil {
		e = append(e, fmt.Errorf("can't receive an underling sql.DB instance: %w", err).Error())
	}

	if err := db.Close(); err != nil {
		e = append(e, fmt.Errorf("can't close database connection: %w", err).Error())
	}

	if err := a.mailer.Close(); err != nil {
		e = append(e, fmt.Errorf("can't close connection to an SMTP server").Error())
	}

	if len(e) > 0 {
		return errors.New(strings.Join(e, "/n"))
	}

	return nil
}

func (a *App) Init() {
	s := a.server

	s.Use(logger.New())

	s.Static("/a", "./assets")

	s.Get("/", a.mainView)
	s.Get("/programme-overview", a.programOverviewView)
	s.Get("/keynote-speakers", a.keynoteSpeakersView)
	s.Get("/requirements", a.requirementsView)
	s.Get("/general-information", a.generalInfoView)
	s.Get("/registration-and-submission", a.registrationView)
	s.Post("/registration-and-submission", a.registerNewParticipant)
	s.Get("/upload", a.uploadView)
	s.Post("/upload", a.uploadArticleOrTezisi)

	admin := s.Group("/admin", basicauth.New(basicauth.Config{
		Users: map[string]string{
			"admin": os.Getenv("ADMIN_PASSWORD"),
		},
	}))
	admin.Get("/", a.adminView)
	admin.Get("/file", a.downloadFile)
	admin.Post("/mailing", a.mailing)

	s.Use(a.notFoundView)
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Err with load env %s", err)
	}
}

func main() {
	app, err := NewApp()
	if err != nil {
		log.Fatal(err)
	}

	app.Init()

	stop := make(chan os.Signal, 1)

	signal.Notify(stop, os.Interrupt)

	// Todo: wait for gorutine to see the output
	go func() {
		<-stop

		log.Println("Received an interrupt signal, shutdown")

		// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		// defer cancel()

		if err := app.Shutdown(context.TODO()); err != nil {
			// Do recovery???
			// app.logger.Errorf("Server shutdown failed: %w", err)
			log.Printf("Server shutdown failed: %v", err)
		}

		log.Println("Seccessfully shutdowned")

	}()

	app.Run()
}
