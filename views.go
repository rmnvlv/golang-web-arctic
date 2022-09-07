package main

import (
	"embed"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

//go:embed views
var ViewsFS embed.FS

func (a *App) mainView(c *fiber.Ctx) error {
	c.Bind(fiber.Map{
		"Header":  true,
		"Content": IndexPageContent,
	})
	return c.Render("index", fiber.Map{})
}

func (a *App) programOverviewView(c *fiber.Ctx) error {
	c.Bind(fiber.Map{
		"Title": "Programme Overview",
	})
	return c.Render("programm-overview", fiber.Map{})
}

func (a *App) keynoteSpeakersView(c *fiber.Ctx) error {
	c.Bind(fiber.Map{
		"Title":   "Keynote Speakers",
		"Content": "Key speakers to be determined later.",
	})
	return c.Render("basic", fiber.Map{})
}

func (a *App) requirementsView(c *fiber.Ctx) error {
	c.Bind(fiber.Map{
		"Title":   "Requirements",
		"Content": "Article template will be posted later.",
	})
	return c.Render("basic", fiber.Map{})
}

func (a *App) generalInfoView(c *fiber.Ctx) error {
	c.Bind(fiber.Map{
		"Title": "Genral Information",
	})
	return c.Render("general-information", fiber.Map{})
}

func (a *App) registrationView(c *fiber.Ctx) error {
	c.Bind(fiber.Map{
		"Title": "Register",
	})
	return c.Render("registration", fiber.Map{})
}

func (a *App) adminView(c *fiber.Ctx) error {
	var participants []Participant
	if err := a.db.Find(&participants).Error; err != nil {
		a.log.Error(err)
		c.Bind(fiber.Map{"Errors": map[string]string{"getParticipants": "Can't fetch participants"}})
	} else {
		c.Bind(fiber.Map{"Users": participants})
	}

	return c.Render("admin", fiber.Map{})
}

func (a *App) preRegistration(c *fiber.Ctx) error {
	c.Bind(fiber.Map{
		"Title": "Verify email",
	})

	return c.Render("preregistration", fiber.Map{})
}

func (a *App) notFoundView(c *fiber.Ctx) error {
	c.Bind(fiber.Map{
		"Title":   "Page Not Found",
		"Content": fmt.Sprintf("Not found - %s", c.OriginalURL()),
	})
	c.Status(fiber.StatusNotFound)
	return c.Render("basic", fiber.Map{})
}
