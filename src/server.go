package main

import (
	scraper_control "genote-watcher/scraper-control"
	"log"
	"net/http"
)

var command_channel chan scraper_control.ScraperCommandType

func registerRoutes() {
	http.Handle("/", http.FileServer(http.Dir("./client/dist")))

	http.HandleFunc("/api/scraper/stop", func(w http.ResponseWriter, r *http.Request) {
		command_channel <- scraper_control.Stop
		log.Printf("Hello API was called")
		w.Write([]byte("Hello, World!"))
	})
}

func StartServer(command_c chan scraper_control.ScraperCommandType) {
	command_channel = command_c
	registerRoutes()
	log.Printf("Server is running on port 4000")
	http.ListenAndServe(":4000", nil)
}
