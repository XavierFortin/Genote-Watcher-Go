package main

import (
	"log"
	"runtime/debug"
	"time"

	"genote-watcher/model"
	scraper_control "genote-watcher/scraper-control"
	"genote-watcher/utils"
)

var config *model.Config
var BuildMode string
var ScraperCommandChannel chan scraper_control.ScraperCommandType

func main() {
	defer func() {
		if r := recover(); r != nil {
			stackTrace := string(debug.Stack())
			log.Println(stackTrace)

			if BuildMode == "prod" {
				utils.NotifyOnCrash(config.DiscordWebhook)
			}
		}
	}()

	config = utils.MustGetConfig()

	if config.TimeInterval == 0 {
		StartGenoteScraping(config)
	} else {
		go func() {
			currentTimeInterval := config.TimeInterval
			ticker := time.Tick(currentTimeInterval)

			for {
				select {
				case <-ticker:
					log.Println("Fake getting Genote Scraping")
					currentTimeInterval = currentTimeInterval - 1*time.Second
					log.Printf("Current Time Interval: %s\n", currentTimeInterval)
					//StartGenoteScraping(config)
				case command := <-ScraperCommandChannel:
					switch command {
					case scraper_control.Restart:
						log.Println("Restarting Genote Scraping")
					case scraper_control.ForceStart:
						log.Println("Force Starting Genote Scraping")
					}

				}
			}
		}()
	}

	StartServer()
}
