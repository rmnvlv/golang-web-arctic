package main

type Participant struct {
	CreatedAt    string
	CodeUUID     string `gorm:"primaryKey"`
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

	Article string
	Tizis   string
}
