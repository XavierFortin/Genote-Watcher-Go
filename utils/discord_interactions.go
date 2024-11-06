package utils

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
)

func NotifyUser(url, courseCode string) {
	data := []byte(fmt.Sprintf(`
	{
    "content": "@everyone",
    "embeds": [
      {
        "description": "Nouvelle note en %s est disponible sur Genote!",
        "color": 100425
      }
    ]
	}
	`, courseCode))

	contentType := "application/json"
	r, err := http.Post(url, contentType, bytes.NewBuffer(data))
	HandleFatalError(err)

	defer r.Body.Close()

	if r.StatusCode == 204 {
		log.Printf("Notification sent successfully")
	} else {
		log.Println("Failed to send notification")
	}
}

func NotifyOnCrash(url string) {
	data := []byte(`
	{
		"embeds": [
			{
				"description": "Genote Watcher a crashé!",
				"color": 16711680
			}
		]
	}
	`)

	contentType := "application/json"
	r, err := http.Post(url, contentType, bytes.NewBuffer(data))
	HandleFatalError(err)
	defer r.Body.Close()

	if r.StatusCode == 204 {
		log.Printf("Crash notification sent successfully")
	} else {
		log.Println("Failed to send crash notification")
	}
}
