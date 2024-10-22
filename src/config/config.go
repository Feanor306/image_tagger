package config

import (
	"os"

	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
	SrvPort    string
	DbUser     string
	DbPassword string
	DbHost     string
	DbPort     string
	DbName     string
}

var conf *Config

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
