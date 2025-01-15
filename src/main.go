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
			ticker := time.NewTicker(config.TimeInterval)

			for {
				select {
				case <-ticker.C:
					log.Println("Fake getting Genote Scraping")
					//StartGenoteScraping(config)
				case command := <-ScraperCommandChannel:
					switch command {
					case scraper_control.Restart:
						log.Println("Restarting Genote Scraping")
						ticker.Reset(config.TimeInterval)

					case scraper_control.ForceStart:
						log.Println("Force Starting Genote Scraping")
						StartGenoteScraping(config)

					case scraper_control.Stop:
						log.Println("Stopping Genote Scraping")
						ticker.Stop()

					default:
					}
				}
			}
		}()
	}

	StartServer()
}
