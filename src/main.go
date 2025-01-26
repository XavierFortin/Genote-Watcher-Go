package main

import (
	"flag"
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
				config := config.MustGetConfig()

				utils.NotifyOnCrash(config.DiscordWebhook)
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

	var port string
	flag.StringVar(&port, "port", "4000", "port to run the server on")

	flag.Parse()

	StartServer(&scraper, port)
}
