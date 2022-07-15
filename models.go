package main

import "gorm.io/gorm"

type Participant struct {
	gorm.Model
	Surname             string `json:"surname"`
	Name                string `json:"name"`
	Organization        string `json:"organization"`
	Position            string `json:"position"`
	Phone               string `json:"phone"`
	Email               string `json:"email"`
	Type                string `json:"type"`         // Speaker/Publication/Listener
	Presentation        string `json:"presentation"` // Plenary session/Section 1/Section 2/Section 3/Section 4
	TitleOfPresentation string `json:"titleofpresentation"`
}

type Error struct {
	Phone   string
	Email   string
	Name    string
	Surname string
}
