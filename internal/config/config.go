package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

// Config is a config
type Config struct {
	HttpAddr           string
	SecretKey          string
	MongoUrl           string
	RedisUrl           string
	HealthcheckMessage string
	ClientOrigin       string
}

// ReadEnv Read reads config from environment.
func ReadEnv() Config {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}

	var config Config
	if httpAddr, exists := os.LookupEnv("HTTP_ADDR"); exists {
		config.HttpAddr = httpAddr
	}
	if secretKey, exists := os.LookupEnv("SECRET_KEY"); exists {
		config.SecretKey = secretKey
	}
	if mongoUrl, exists := os.LookupEnv("MONGODB_URL"); exists {
		config.MongoUrl = mongoUrl
	}
	if redisUrl, exists := os.LookupEnv("REDIS_URL"); exists {
		config.RedisUrl = redisUrl
	}

	if healthcheckMessage, exists := os.LookupEnv("HEALTHCHECK_MESSAGE"); exists {
		config.HealthcheckMessage = healthcheckMessage
	} else {
		config.HealthcheckMessage = "ok"
	}

	if clientOrigin, exists := os.LookupEnv("CLIENT_ORIGIN"); exists {
		config.ClientOrigin = clientOrigin
	}
	return config
}
