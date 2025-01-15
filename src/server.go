package main

import (
	"log"
	"net/http"
)

func registerRoutes() {
	http.Handle("/", http.FileServer(http.Dir("./client/dist")))

	http.HandleFunc("/api/hello", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Hello API was called")
		w.Write([]byte("Hello, World!"))
	})
}

func StartServer() {
	registerRoutes()
	log.Printf("Server is running on port 4000")
	http.ListenAndServe(":4000", nil)
}
