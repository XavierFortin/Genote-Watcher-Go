package utils

import (
	"encoding/json"
	"genote-watcher/model"
	"log"
	"math/rand"
	"os"

	"github.com/joho/godotenv"
)

func GetUserAgents() []string {
	return []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:130.0) Gecko/20100101 Firefox/130.0",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 Safari/537.36 Edg/128.0.0.0",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/127.0.0.0 Safari/537.36 OPR/113.0.0.0",
	}
}

func GetRandomUserAgent() string {
	userAgents := GetUserAgents()
	return userAgents[rand.Intn(len(userAgents))]
}

func GetEnvVariable(key string) string {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

func CreateResultFile(data []model.CourseRow) {
	r, _ := json.Marshal(data)

	err := os.WriteFile("result.json", r, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func ReadResultFile() []model.CourseRow {
	file, err := os.ReadFile("result.json")
	if err != nil {
		log.Fatal(err)
	}

	var data []model.CourseRow
	err = json.Unmarshal(file, &data)
	if err != nil {
		log.Fatal(err)
	}

	return data
}