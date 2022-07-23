package main

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

type link struct {
	Link string
	Text string
}

func links(source []string) []link {
	links := make([]link, len(source))
	for i, s := range source {
		links[i].Link = "/" + strings.ReplaceAll(strings.ToLower(s), " ", "-")
		links[i].Text = s
	}
	return links
}

var Links = links([]string{"Programme Overview", "Keynote Speakers", "Registration and submission", "Requirements", "General information"})

var IndexPage = fiber.Map{
	// "Title": "",
	"Content": []fiber.Map{
		{
			"Title":     "Welcome to AMTC 2022",
			"Paragraph": "Welcome to the International Conference «Arctic: Marine Transportation Challenges – 2022» (AMTC-2022)! It is the favorite time when the problems of the maritime transport and ecology in the Arctic will be discussed globally at the University campus. The Conference to become as “smart” platform for discussion and exchange of ideas between International organizations, executive bodies, business community, research and educational organizations. Apart from main sessions, after COVID-19 we are plan to the face-to-face Conference, we will be happy to organize an extensive cultural program in Saint-Petersburg with excursions and informal cocktail. We are looking forward to seeing You as a participant!",
		},
		{
			"Title":     "About AMTC-2021",
			"Paragraph": "On October 15, 2021, The International Conference «Arctic: Marine Transportation Challenges – 2021» was held at the Admiral Makarov State University of Maritime and Inland Shipping. Due to the restrictions associated with COVID-19 pandemic, the conference sessions were held in mixed offline and online formats. It was held with the support of the Russian Maritime Register of Shipping, Sovcomflot, Administration of the Northern Sea Route, the Arctic and Antarctic Research Institute, with the participation of the leaders of the World Maritime University and the International Association of Maritime Universities. The topics of the plenary reports included the issues of shipping, shipbuilding,port activities, training for work in the Arctic. The conference participants represented Russia, Norway, Finland, USA, Sweden, Japan and China.",
		},
	},
}

var RegistrationPage = map[string]interface{}{
	"Title": "",
	"Content": fiber.Map{
		"ConferenceSessions": []string{
			"Plenary session",
			"Session 1. Education and professional training for the Arctic shipping industry",
			"Session 2. Arctic shipping: safety, environment, legal regulation",
			"Session 3. Innovative technologies for polar shipping: research & development",
			"Session 4. Development of Sea Ports in the Arctic",
		},
		"ParticipationForm": []string{
			"Speaker",
			"Publication",
			"Speaker & Publication",
			"Listener",
		},
	},
}
