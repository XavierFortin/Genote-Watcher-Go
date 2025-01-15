package main

import (
	"genote-watcher/model"
	"genote-watcher/scrapers"
	"genote-watcher/utils"
	"log"

	"github.com/gocolly/colly/v2"
)

const (
	LOGIN_URL = "https://cas.usherbrooke.ca/login?service=https://www.usherbrooke.ca/genote/public/index.php"
)

func StartGenoteScraping(config *model.Config) {
	c := utils.CreateCollector()

	fieldsData := map[string]string{
		"username": config.Username,
		"password": config.Password,
		"submit":   "",
	}

	c.OnHTML("input[type='hidden']", func(e *colly.HTMLElement) {
		fieldsData[e.Attr("name")] = e.Attr("value")
	})

	c.Visit(LOGIN_URL)

	err := c.Post(LOGIN_URL, fieldsData)
	if err != nil {
		log.Println("Error while logging in: ")
		log.Println(err)
	}

	rows := scrapers.ScrapeCourseRows(c.Clone())

	oldRows := utils.ReadResultFile()
	if oldRows == nil {
		utils.WriteResultFile(rows)
		return
	}

	verifyForChanges(rows, oldRows)

	utils.WriteResultFile(rows)
}

func verifyForChanges(newRows, oldRows []model.CourseRow) {
	diffRows := []string{}

	for index := range newRows {
		if !newRows[index].Equal(&oldRows[index]) {
			diffRows = append(diffRows, newRows[index].CourseCode)
		}
	}

	var changesDetected bool
	for _, courseCode := range diffRows {
		log.Printf("New grade in %s is available on Genote!\n", courseCode)

		if BuildMode == "prod" {
			utils.NotifyUser(config.DiscordWebhook, courseCode)
		}
		changesDetected = true
	}

	if !changesDetected {
		log.Println("No changes detected")
	}
}
