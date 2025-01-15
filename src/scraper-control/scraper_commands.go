package scraper_control

type ScraperCommandType int

const (
	Restart ScraperCommandType = iota
	ForceStart
)
