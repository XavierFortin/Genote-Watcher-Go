package scraper_control

type ScraperCommandType int

const (
	Start ScraperCommandType = iota
	Stop
	Status
	Restart
	ForceStartOnce
)

type Command struct {
	Action ScraperCommandType
}

type Response interface{}

type StatusResponse struct {
	IsRunning bool `json:"isRunning"`
}
