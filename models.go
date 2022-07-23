package main

import "gorm.io/gorm"

type Participant struct {
	gorm.Model
	Id           int `gorm:"primaryKey"`
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

type LoadedFile struct {
	gorm.Model
	Id       int `gorm:"primaryKey"`
	File     string
	FileName string
}

type FormError struct {
	Message string
	Phone   string
	Email   string
	Name    string
	Surname string
	Captcha string
}
