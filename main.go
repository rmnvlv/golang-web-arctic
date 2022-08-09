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
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/gomail.v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var Cfg = make(map[string]string)

type App struct {
	server *fiber.App
	db     *gorm.DB
	mailer gomail.SendCloser
	log    *zap.SugaredLogger
	disk   Disk
	// captcha *HCaptcha
}

func NewApp(log *zap.SugaredLogger) (*App, error) {
	dbURL := Cfg["DATABASE_URL"]

	// TODO: .env DB_URL="/name" + pull docker with my sql
	// db, err := gorm.Open(mysql.Open(dbURL), &gorm.Config{})
	// if err != nil {
	// 	return nil, fmt.Errorf("can't open database: %w", err)
	// }

	db, err := gorm.Open(sqlite.Open(dbURL), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("can't open database")
	}

	log.Infof("Connected to database: %s", dbURL)

	if err := db.AutoMigrate(&Participant{}); err != nil {
		return nil, fmt.Errorf("can't apply migrations to database: %w", err)
	}

	log.Info("Migrations applied")

	smtpPort, err := strconv.Atoi(Cfg["SMTP_PORT"])
	if err != nil {
		return nil, fmt.Errorf("can't convert an SMTP server port to int: %w", err)
	}

	m, err := gomail.NewDialer(Cfg["SMTP_HOST"], smtpPort, Cfg["SMTP_USER"], Cfg["SMTP_PASSWORD"]).Dial()
	if err != nil {
		return nil, fmt.Errorf("can't authenticate to an SMTP server: %w", err)
	}

	log.Infof("Authenticated to SMTP server: %s:%d", Cfg["SMTP_HOST"], smtpPort)

	disk, err := NewOsDisk(os.Getenv("DISK_PATH"))
	if err != nil {
		return nil, fmt.Errorf("can't create disk: %w", err)
	}

	log.Info("Disk initialized at: %s", disk.Path)

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
	}

	app.bootstrap()

	log.Info("Initialization finished")

	return &app, nil
}

func (a *App) Run() {
	// ln, err := net.Listen("unix", "/tmp/arctic.sock")
	// if err != nil {
	// 	a.log.Fatal("Listen error: ", err)
	// }

	// a.server.Listener(ln)

	a.server.Listen(":" + Cfg["HOST"])
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
			"admin": Cfg["ADMIN_PASSWORD"],
		},
	}))
	admin.Get("/", a.adminView)
	admin.Get("/file", a.downloadFile)
	admin.Post("/mailing", a.sendMailing)

	s.Use(a.notFoundView)
}

func CheckEnv() {
	values := []string{"HCAPTCHA_SECRET_KEY", "HCAPTCHA_SECRET_KEY", "YANDEX_OAUTH_TOKEN",
		"SMTP_USER", "SMTP_PASSWORD", "SMTP_HOST", "SMTP_PORT", "ADMIN_PASSWORD",
		"DATABASE_URL", "HOST", "DOMAIN"}

	for _, value := range values {
		path, exists := os.LookupEnv(value)
		if path == "" || !exists {
			log.Fatalf("%s does not exists, please fill .env", value)
		}
		Cfg[value] = path
	}

	log.Print(Cfg)
}

func init() {
	if err := godotenv.Load(); err != nil {
		panic("can't load .env: " + err.Error())
	}

	CheckEnv()
}

func main() {
	lc := zap.NewDevelopmentConfig()
	lc.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000")
	lc.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	l, err := lc.Build()
	if err != nil {
		panic(err)
	}

	log := l.Sugar()

	log.Info("Logger initialized")

	app, err := NewApp(log)
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
