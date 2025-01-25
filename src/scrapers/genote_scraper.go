package scrapers

import (
	"cmp"
	"fmt"
	"genote-watcher/config"
	"genote-watcher/model"
	scraper_control "genote-watcher/scraper-control"
	"genote-watcher/utils"
	"log"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
)

const (
	LOGIN_URL = "https://cas.usherbrooke.ca/login?service=https://www.usherbrooke.ca/genote/public/index.php"
)

type GenoteScraper struct {
	isRunning   bool
	config      config.Config
	ticker      *time.Ticker
	CommandChan chan scraper_control.Command
	ReponseChan chan scraper_control.Response
}

// Creates a new genoteScraper. Environment variables need to exist to create a new genoteScraper
func NewGenoteScraper() GenoteScraper {
	config := config.MustGetConfig()
	return GenoteScraper{
		isRunning:   true,
		config:      config,
		ticker:      nil,
		CommandChan: make(chan scraper_control.Command),
		ReponseChan: make(chan scraper_control.Response),
	}
}

func (gs *GenoteScraper) Start() {
	gs.isRunning = true
	if gs.config.TimeInterval == 0 {
		gs.ScrapeOnce()
	} else {
		go func() {
			gs.ticker = time.NewTicker(gs.config.TimeInterval)
			for {
				select {
				case <-gs.ticker.C:
					gs.ScrapeOnce()

				case command := <-gs.CommandChan:
					gs.handleCommand(command)
				}
			}
		}()
	}
}

func (gs *GenoteScraper) handleCommand(command scraper_control.Command) {
	switch command.Action {
	case scraper_control.Start:
		gs.ticker.Reset(gs.config.TimeInterval)
		gs.isRunning = true

	case scraper_control.Stop:
		gs.ticker.Stop()
		gs.isRunning = false

	case scraper_control.Status:
		gs.ReponseChan <- scraper_control.StatusResponse{IsRunning: gs.isRunning, Interval: gs.config.TimeInterval.String()}

	case scraper_control.Restart:
		log.Printf("TimeInterval: %s\n", gs.config.TimeInterval)
		gs.ticker.Reset(gs.config.TimeInterval)
		gs.isRunning = true
		gs.ScrapeOnce()

	case scraper_control.ForceStartOnce:
		gs.ScrapeOnce()

	case scraper_control.ChangeInterval:
		duration, err := time.ParseDuration(command.Data.(string))

		utils.HandleLogError(err)
		gs.config.SetTimeInterval(duration)

		log.Printf("New interval: %s\n", gs.config.TimeInterval)
	default:
	}
}

func (gs *GenoteScraper) ScrapeOnce() {
	c := utils.CreateCollector()

	fieldsData := map[string]string{
		"username": gs.config.Username,
		"password": gs.config.Password,
		"submit":   "",
	}

	c.OnHTML("input[type='hidden']", func(e *colly.HTMLElement) {
		fieldsData[e.Attr("name")] = e.Attr("value")
	})

	c.Visit(LOGIN_URL)

	err := c.Post(LOGIN_URL, fieldsData)
	if err != nil {
		log.Println("Error while logging in: ")
		log.Println(err)
	}

	rows := scrapeCourseRows(c.Clone())

	oldRows := utils.ReadResultFile()
	if oldRows == nil {
		utils.WriteResultFile(rows)
		return
	}

	gs.verifyForChanges(rows, oldRows)

	utils.WriteResultFile(rows)

}

// ScrapeCourses scrapes the courses from the genote website
// and returns the amount of empty notes found
func scrapeCoursesEmptyNotes(c *colly.Collector, url string) int {
	if url == "" {
		return 0
	}

	emptyNotes := 0

	c.OnHTML("table.zebra tbody td:nth-child(3)", func(e *colly.HTMLElement) {
		if strings.Contains(e.Text, "-- /") {
			emptyNotes++
		}
	})

	err := c.Visit(fmt.Sprintf("https://www.usherbrooke.ca/genote/application/etudiant/%s", url))
	utils.HandleLogError(err)

	c.Wait()

	return emptyNotes
}

func scrapeCourseRows(c *colly.Collector) []model.CourseRow {
	rows := []model.CourseRow{}
	c.OnHTML("table:nth-child(4) tbody", func(e *colly.HTMLElement) {
		cr := model.CourseRow{}
		e.ForEach("tr", func(_ int, el *colly.HTMLElement) {
			cr.CourseName = el.DOM.Find("td:nth-child(1)").Text()
			splitName := strings.Split(cr.CourseName, " ")
			courseCode := splitName[len(splitName)-2]

			cr.CourseCode = courseCode[1:]
			cr.EvaluationAmount, _ = strconv.Atoi(el.DOM.Find("td:nth-child(5)").Text())
			cr.CourseLink = el.DOM.Find("td:nth-child(6) a").AttrOr("href", "")
			cr.EmptyNoteAmount = scrapeCoursesEmptyNotes(c.Clone(), cr.CourseLink)

			rows = append(rows, cr)
		})
	})

	err := c.Visit("https://www.usherbrooke.ca/genote/application/etudiant/cours.php")
	utils.HandleLogError(err)

	c.Wait()

	// Sort rows by course code
	slices.SortFunc(rows, func(i, j model.CourseRow) int {
		return cmp.Compare(i.CourseName, j.CourseName)
	})

	return rows
}

func (gs *GenoteScraper) verifyForChanges(newRows, oldRows []model.CourseRow) {
	diffRows := []string{}

	for index := range newRows {
		if !newRows[index].Equal(&oldRows[index]) {
			diffRows = append(diffRows, newRows[index].CourseCode)
		}
	}

	var changesDetected bool
	for _, courseCode := range diffRows {
		log.Printf("New grade in %s is available on Genote!\n", courseCode)

		if utils.BuildMode == "prod" {
			utils.NotifyUser(gs.config.DiscordWebhook, courseCode)
		}
		changesDetected = true
	}

	if !changesDetected {
		log.Println("No changes detected")
	}
}
