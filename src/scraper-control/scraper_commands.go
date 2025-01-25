package scraper_control

type ScraperCommandType int

const (
	Start ScraperCommandType = iota
	Stop
	Status
	Restart
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
