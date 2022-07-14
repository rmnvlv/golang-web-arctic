package main

import "gorm.io/gorm"

type Participant struct {
	gorm.Model
	Surname      string `json:"surname"`
	Name         string `json:"name"`
	Organizacion string `json:"organizacion"`
	Position     string `json:"position"`
	Phone        string `json:"phone"`
	Email        string `json:"email"`
	Type         string `json:"type"`     // Speaker/Publication/Listener
	Planning     string `json:"planning"` // Plenary session/Section 1/Section 2/Section 3/Section 4
	Title        string `json:"title"`
}
