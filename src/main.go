package main

import (
	"io"
	"log"
	"os"
	"runtime/debug"

	scraper_control "genote-watcher/scraper-control"
	"genote-watcher/scrapers"
	"genote-watcher/utils"
)

func main() {
	var command_channel chan scraper_control.ScraperCommandType = make(chan scraper_control.ScraperCommandType)
	var scraper = scrapers.NewGenoteScraper(command_channel)

	logFile, err := os.OpenFile("log.txt", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	defer logFile.Close()
	if err != nil {
		panic(err)
	}

	chanWriter := utils.NewChannelWriter(150)
	defer chanWriter.Close()

	mw := io.MultiWriter(os.Stdout, logFile, chanWriter)
	log.SetOutput(mw)

	defer func() {
		if r := recover(); r != nil {
			stackTrace := string(debug.Stack())
			log.Println(stackTrace)

			if utils.BuildMode == "prod" {
				utils.NotifyOnCrash(utils.MustGetConfig().DiscordWebhook)
			}
		}
	}()

	scraper.Start()

	StartServer(command_channel)
}
