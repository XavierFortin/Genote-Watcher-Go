package utils

import (
	"encoding/json"
	"errors"
	"genote-watcher/model"
	"log"
	"math/rand"
	"net/http/cookiejar"
	"os"

	"github.com/gocolly/colly/v2"
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

func CreateCollector() *colly.Collector {
	c := colly.NewCollector(
		colly.UserAgent(GetRandomUserAgent()),
	)

	jar, _ := cookiejar.New(nil)
	c.SetCookieJar(jar)

	return c
}

func WriteResultFile(data []model.CourseRow) {
	r, _ := json.Marshal(data)

	err := os.WriteFile("result.json", r, 0644)
	HandleFatalError(err)
}

func ReadResultFile() []model.CourseRow {

	if _, err := os.Stat("result.json"); errors.Is(err, os.ErrNotExist) {
		os.Create("result.json")
		return nil
	}

	file, err := os.ReadFile("result.json")

	HandleFatalError(err)

	var data []model.CourseRow
	err = json.Unmarshal(file, &data)
	HandleFatalError(err)

	if len(data) == 0 {
		return nil
	}

	return data
}

func HandleLogError(err error) {
	if err != nil {
		log.Println(err)
	}
}

func HandleFatalError(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
