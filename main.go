package main

import (
	"fmt"
	"time"

	"genote-watcher/model"
	"genote-watcher/scrapers"
	"genote-watcher/utils"

	"github.com/gocolly/colly/v2"
)

const (
	LOGIN_URL = "https://cas.usherbrooke.ca/login?service=https://www.usherbrooke.ca/genote/public/index.php"
)

var config *utils.Config

func login(c *colly.Collector) {

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
		fmt.Println("Error while logging in: ")
		fmt.Println(err)
	}
}

func notifyForChanges(newRows, oldRows []model.CourseRow) {
	diffRows := []string{}

	for index := range newRows {
		if !newRows[index].Equal(&oldRows[index]) {
			diffRows = append(diffRows, newRows[index].CourseCode)
		}
	}

	formattedDate := time.Now().Format("2006/01/02 15:04:05")
	var changesDetected bool
	for _, courseCode := range diffRows {
		fmt.Printf("[%s] Nouvelle note en %s est disponible sur Genote!\n", formattedDate, courseCode)
		utils.NotifyUser(config.DiscordWebhook, courseCode)
		changesDetected = true
	}

	if !changesDetected {
		fmt.Printf("[%s] Aucun changement détecté\n", formattedDate)
	}
}

func startGenoteScraping() {
	c := utils.CreateCollector()

	login(c)

	rows := scrapers.ScrapeCourseRows(c.Clone())

	oldRows := utils.ReadResultFile()
	if oldRows == nil {
		utils.WriteResultFile(rows)
		return
	}

	notifyForChanges(rows, oldRows)

	utils.WriteResultFile(rows)
}

func main() {
	config = utils.MustGetConfig()

	if config.TimeInterval == 0 {
		startGenoteScraping()
	} else {
		for range time.Tick(config.TimeInterval) {
			startGenoteScraping()
		}
	}
}
