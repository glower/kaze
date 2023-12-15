package config

import "os"

type Config struct {
	DB              string
	OpenMeteoAPIKey string
	IsDebug         bool
}

func NewConfig() *Config {
	return &Config{
		DB:              os.Getenv("APP_DB"),
		OpenMeteoAPIKey: os.Getenv("OPEN_METEO_API_KEY"),
		IsDebug:         os.Getenv("DEBUG") != "",
	}
}
