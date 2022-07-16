package main

import "gorm.io/gorm"

type Participant struct {
	gorm.Model

	Surname      string
	Name         string
	Organization string
	Position     string
	Phone        string
	Email        string

	// Speaker|Publication|Listener
	PresentationForm string

	// Plenary session, etc.
	PresentationSection string

	PresentationTitle string
}

type FormError struct {
	Message string
	Phone   string
	Email   string
	Name    string
	Surname string
	Captcha string
}
