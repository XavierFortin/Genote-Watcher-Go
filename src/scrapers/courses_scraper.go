package scrapers

import (
	"cmp"
	"fmt"
	"genote-watcher/model"
	"genote-watcher/utils"
	"slices"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
)

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

func ScrapeCourseRows(c *colly.Collector) []model.CourseRow {
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
