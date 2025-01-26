package scrapers

import (
	"cmp"
	"fmt"
	"genote-watcher/config"
	"genote-watcher/model"
	"genote-watcher/scraper_commands"
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
	isRunning    bool
	isConfigured bool
	config       config.Config
	ticker       *time.Ticker
}

// Creates a new genoteScraper. Environment variables need to exist to create a new genoteScraper
func NewGenoteScraper() GenoteScraper {
	config := config.MustGetConfig()
	return GenoteScraper{
		isRunning:    false,
		isConfigured: true,
		config:       config,
		ticker:       nil,
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
				}
			}
		}()
	}
}

func (gs *GenoteScraper) Stop() {
	gs.ticker.Stop()
	gs.isRunning = false
}

func (gs *GenoteScraper) Resume() {
	if !gs.isRunning {
		gs.ticker.Reset(gs.config.TimeInterval)
		gs.isRunning = true
	}
}

func (gs *GenoteScraper) GetStatus() scraper_commands.StatusResponse {
	return scraper_commands.StatusResponse{IsRunning: gs.isRunning, Interval: gs.config.TimeInterval.String()}
}

func (gs *GenoteScraper) SetInterval(duration time.Duration) {
	gs.config.SetTimeInterval(duration)
	gs.ticker.Reset(gs.config.TimeInterval)
	if !gs.isRunning {
		gs.ticker.Stop()
		gs.isRunning = false
	}

	log.Printf("New interval: %s\n", gs.config.TimeInterval)
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
