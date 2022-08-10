package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"gopkg.in/gomail.v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// var Cfg = make(map[string]string)

type App struct {
	server *fiber.App
	db     *gorm.DB
	mailer gomail.SendCloser
	log    *zap.SugaredLogger
	disk   Disk
	config Config
}

func NewApp(config Config, log *zap.SugaredLogger) (*App, error) {

	// db, err := gorm.Open(mysql.Open(config.DatabaseURL), &gorm.Config{})
	db, err := gorm.Open(sqlite.Open(config.DatabaseURL), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("can't open database: %w", err)
	}
	log.Infof("Connected to database: %s", config.DatabaseURL)

	if err := db.AutoMigrate(&Participant{}); err != nil {
		return nil, fmt.Errorf("can't apply migrations to database: %w", err)
	}
	log.Info("Migrations applied")

	m, err := gomail.NewDialer(config.SMTP.Host, config.SMTP.Port, config.SMTP.User, config.SMTP.Password).Dial()
	if err != nil {
		return nil, fmt.Errorf("can't authenticate to an SMTP server: %w", err)
	}
	log.Infof("Authenticated to SMTP server: %s:%d", config.SMTP.Host, config.SMTP.Port)

	disk, err := NewOsDisk(config.DiskPath)
	if err != nil {
		return nil, fmt.Errorf("can't create disk: %w", err)
	}
	log.Infof("Disk initialized at: %s", disk.Path)

	server := fiber.New(fiber.Config{
		Views:       html.New("./views", ".html"),
		ViewsLayout: "main",
	})

	app := App{
		server: server,
		db:     db,
		mailer: m,
		log:    log,
		disk:   disk,
		config: config,
	}

	app.bootstrap()

	log.Info("Initialization finished")

	return &app, nil
}

func (a *App) Run() {
	if a.config.HTTPAddressUnix != "" {
		ln, err := net.Listen("unix", a.config.HTTPAddressUnix)
		if err != nil {
			a.log.Fatal("Listen error: ", err)
		}
		a.server.Listener(ln)
	} else if a.config.HTTPAddress != "" {
		a.server.Listen(a.config.HTTPAddress)
	} else {
		a.server.Listen(":" + os.Getenv("PORT"))
	}
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

func (a *App) bootstrap() {
	s := a.server

	s.Use(
		logger.New(),
		// NewLoggerMiddleware(Config{Logger: a.log.Desugar(), Next: nil}),
	)

	s.Static("/a", "./assets")

	s.Get("/", a.mainView)
	s.Get("/programme-overview", a.programOverviewView)
	s.Get("/keynote-speakers", a.keynoteSpeakersView)
	s.Get("/requirements", a.requirementsView)
	s.Get("/general-information", a.generalInfoView)
	s.Get("/registration-and-submission", a.registrationView)
	s.Post("/registration-and-submission", a.registerNewParticipant)
	s.Get("/upload/:type", a.uploadView)
	s.Post("/upload/:type", a.uploadFile)

	admin := s.Group("/admin", basicauth.New(basicauth.Config{
		Users: map[string]string{
			"admin": a.config.AdminPassword,
		},
	}))
	admin.Get("/", a.adminView)
	admin.Get("/file", a.downloadExcel)
	admin.Post("/mailing", a.sendMailing)
	admin.Get("/download/:file", a.downloadFiles)

	s.Use(a.notFoundView)
}

func init() {
	if err := godotenv.Load(); err != nil {
		panic("can't load .env: " + err.Error())
	}
}

func main() {
	log, err := NewLogger()
	if err != nil {
		panic(err)
	}
	log.Info("Logger initialized")

	config, err := NewConfig()
	if err != nil {
		log.Fatalf("Config not loaded: %w", err)
	}

	log.Infof("Config loaded:\n%s", config.String())

	app, err := NewApp(config, log)
	if err != nil {
		log.Fatal(err)
	}

	stop := make(chan os.Signal, 1)

	signal.Notify(stop, os.Interrupt)

	// Todo: wait for gorutine to see the output
	go func() {
		<-stop

		log.Info("Received an interrupt signal, shutdown")

		// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		// defer cancel()

		if err := app.Shutdown(context.TODO()); err != nil {
			// Do recovery???
			// app.logger.Errorf("Server shutdown failed: %w", err)
			log.Errorf("Server shutdown failed: %v", err)
		}

		log.Info("Seccessfully shutdowned")

	}()

	app.Run()
}
