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
	"github.com/spf13/cobra"
	"gopkg.in/gomail.v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func init() {
	if err := godotenv.Load(); err != nil {
		panic("can't load .env: " + err.Error())
	}
}

type App struct {
	server *fiber.App
	db     *gorm.DB
	mailer gomail.SendCloser
	log    *Logger
	disk   Disk
	config *Config
}

func (a *App) Init(config *Config, log *Logger) error {
	// db, err := gorm.Open(mysql.Open(config.DatabaseURL), &gorm.Config{})
	db, err := gorm.Open(sqlite.Open(config.DatabaseURL), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("can't open database: %w", err)
	}
	log.Infof("Connected to database: %s", config.DatabaseURL)

	if err := db.AutoMigrate(&Participant{}); err != nil {
		return fmt.Errorf("can't apply migrations to database: %w", err)
	}
	log.Info("Migrations applied")

	m, err := gomail.NewDialer(config.SMTP.Host, config.SMTP.Port, config.SMTP.User, config.SMTP.Password).Dial()
	if err != nil {
		return fmt.Errorf("can't authenticate to an SMTP server: %w", err)
	}
	log.Infof("Authenticated to SMTP server: %s:%d", config.SMTP.Host, config.SMTP.Port)

	disk, err := NewOsDisk(config.DiskPath)
	if err != nil {
		return fmt.Errorf("can't init disk at '%s': %w", config.DiskPath, err)
	}
	log.Infof("Disk initialized at: %s", disk.Path)

	server := fiber.New(fiber.Config{
		Views:       html.New("./views", ".html"),
		ViewsLayout: "main",
	})

	a.server = server
	a.db = db
	a.mailer = m
	a.log = log
	a.disk = disk
	a.config = config

	a.registerRoutes()

	return nil
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

func (a *App) registerRoutes() {
	s := a.server

	s.Use(
		logger.New(),
		// NewLoggerMiddleware(Config{Logger: a.log.Desugar(), Next: nil}),
	)

	s.Static("/a", "./assets")

	s.Use(func(c *fiber.Ctx) error {
		c.Bind(fiber.Map{
			"Links": links([]string{"Programme Overview", "Keynote Speakers", "Registration and submission", "Requirements", "General information"}),
		})
		return c.Next()
	})

	s.Get("/", a.mainView)
	s.Get("/programme-overview", a.programOverviewView)
	s.Get("/keynote-speakers", a.keynoteSpeakersView)
	s.Get("/requirements", a.requirementsView)
	s.Get("/general-information", a.generalInfoView)
	s.Get("/registration-and-submission", a.registrationView)
	s.Post("/registration-and-submission", a.registerNewParticipant)
	s.Get("/upload/:type", a.uploadView)
	s.Post("/upload/:type", a.uploadFile)

	admin := s.Group("/admin",
		basicauth.New(
			basicauth.Config{
				Users: map[string]string{
					"admin": a.config.AdminPassword,
				},
			},
		),
		func(c *fiber.Ctx) error {
			c.Bind(fiber.Map{
				"Title": "Admin",
			})
			return c.Next()
		},
	)
	admin.Get("/", a.adminView)
	admin.Post("/mailing", a.sendNewsletter)
	admin.Get("/download/:file", a.downloadFiles)

	s.Use(a.notFoundView)
}

func main() {
	logger := new(Logger)
	config := new(Config)

	serveCmd := &cobra.Command{
		Use: "serve",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if err := logger.Init(); err != nil {
				return err
			}

			if err := config.LoadEnv(); err != nil {
				return err
			}

			fmt.Println(config)

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			app := new(App)

			if err := app.Init(config, logger); err != nil {
				return fmt.Errorf("init app: %w", err)
			}

			stop := make(chan os.Signal, 1)

			signal.Notify(stop, os.Interrupt)

			// Todo: wait for gorutine to see the output
			go func() {
				<-stop

				logger.Info("Received an interrupt signal, shutdown")

				// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				// defer cancel()

				if err := app.Shutdown(context.TODO()); err != nil {
					// Do recovery???
					logger.Errorf("Application shutdown failed: %v", err.Error())
				}

				logger.Info("Seccessfully stoped application")

			}()

			app.Run()

			return nil
		},
	}

	serveCmd.Flags().StringVar(&config.HTTPAddressUnix, "unix", "", "")
	serveCmd.Flags().StringVar(&config.HTTPAddress, "http", "", "")
	serveCmd.MarkFlagsMutuallyExclusive("unix", "http")

	command := &cobra.Command{
		Use:     "amtc",
		Version: "0.0.0",
	}
	command.PersistentFlags().StringVar(&config.DatabaseURL, "db-url", "", "")
	command.PersistentFlags().StringVar(&config.DiskPath, "disk-path", "", "")
	command.MarkPersistentFlagRequired("db-url")
	command.MarkPersistentFlagRequired("disk-path")
	command.AddCommand(serveCmd)

	command.Execute()
}
