package config

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"sync"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Username       string        `env:"GENOTE_USER" required:"true"`
	Password       string        `env:"GENOTE_PASSWORD" required:"true"`
	DiscordWebhook string        `env:"DISCORD_WEBHOOK" required:"true"`
	TimeInterval   time.Duration `env:"TIME_INTERVAL" required:"false" default:"0"`
}

func (c *Config) SetTimeInterval(duration time.Duration) {
	c.TimeInterval = duration
}

var (
	instance *Config
	once     sync.Once
	loadErr  error
)

type MissingConfigError struct{}

func (e *MissingConfigError) Error() string {
	return "missing configuration"
}

func MustGetConfig() (Config, error) {
	once.Do(func() {
		instance, loadErr = loadEnvVariables()
	})

	if loadErr != nil {
		return Config{}, &MissingConfigError{}
	}

	return *instance, nil
}

func loadEnvVariables() (*Config, error) {
	config := &Config{}

	godotenv.Load()

	t := reflect.TypeOf(config).Elem()
	val := reflect.ValueOf(config).Elem()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := val.Field(i)

		key := field.Tag.Get("env")
		required := field.Tag.Get("required") == "true"
		defaultValue := field.Tag.Get("default")

		envValue := os.Getenv(key)

		if required && envValue == "" {
			if val.Field(i).String() == "" {
				return nil, fmt.Errorf("missing required environment variable %s", key)
			}
		}

		if envValue == "" && defaultValue != "" {
			envValue = defaultValue
		}

		if err := setField(value, envValue); err != nil {
			return nil, fmt.Errorf("failed to set field %s: %v", key, err)
		}
	}

	return config, nil
}

func setField(field reflect.Value, value string) error {

	switch field.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Int:
		i, err := strconv.Atoi(value)
		if err != nil {
			return err
		}
		field.SetInt(int64(i))
	case reflect.TypeOf(time.Duration(0)).Kind():
		d, err := time.ParseDuration(value)
		if err != nil {
			return err
		}
		field.SetInt(int64(d))
	}
	return nil
}
