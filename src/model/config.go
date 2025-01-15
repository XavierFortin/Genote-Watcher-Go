package model

import "time"

type Config struct {
	Username       string        `env:"GENOTE_USER" required:"true"`
	Password       string        `env:"GENOTE_PASSWORD" required:"true"`
	DiscordWebhook string        `env:"DISCORD_WEBHOOK" required:"true"`
	TimeInterval   time.Duration `env:"TIME_INTERVAL" required:"false" default:"0"`
}

func (c *Config) SetTimeInterval(duration time.Duration) {
	c.TimeInterval = duration
}
