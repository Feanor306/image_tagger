package config

import (
	"os"

	_ "github.com/joho/godotenv/autoload"
)

// Config contains environment variables
type Config struct {
	SrvPort    string
	DbUser     string
	DbPassword string
	DbHost     string
	DbPort     string
	DbName     string
}

// conf is the local singleton instance of Config
var conf *Config

// GetConfig will return local config or create it
// with all the relevant env variables
func GetConfig() *Config {
	if conf != nil {
		return conf
	}

	conf = &Config{
		SrvPort:    os.Getenv("TAGGER_SRV_PORT"),
		DbUser:     os.Getenv("TAGGER_DB_USER"),
		DbPassword: os.Getenv("TAGGER_DB_PASSWORD"),
		DbName:     os.Getenv("TAGGER_DB_NAME"),
		DbHost:     os.Getenv("TAGGER_DB_HOST"),
		DbPort:     os.Getenv("TAGGER_DB_PORT"),
	}

	return conf
}
