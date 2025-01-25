package main

import (
	"io"
	"log"
	"os"
	"runtime/debug"

	"genote-watcher/config"
	"genote-watcher/scrapers"
	"genote-watcher/utils"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			stackTrace := string(debug.Stack())
			log.Println(stackTrace)

			if utils.BuildMode == "prod" {
				utils.NotifyOnCrash(config.MustGetConfig().DiscordWebhook)
			}
		}
	}()

	logFile, err := os.OpenFile("log.txt", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	defer logFile.Close()
	if err != nil {
		panic(err)
	}

	mw := io.MultiWriter(os.Stdout, logFile, &WebSocketLogger{})
	log.SetOutput(mw)

	var scraper = scrapers.NewGenoteScraper()
	scraper.Start()

	StartServer(scraper.CommandChan, scraper.ReponseChan)
}
