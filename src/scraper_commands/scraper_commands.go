package scraper_commands

type ScraperCommandType int

const (
	Resume ScraperCommandType = iota
	Stop
	Status
	ForceStartOnce
	ChangeInterval
)

type Command struct {
	Action ScraperCommandType
	Data   interface{}
}

type Response interface{}

type StatusResponse struct {
	IsRunning bool   `json:"isRunning"`
	Interval  string `json:"interval"`
}
