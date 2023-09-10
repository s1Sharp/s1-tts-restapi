package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/s1Sharp/s1-tts-restapi/internal/logger"
)

var (
	log = logger.GetLogger()
)

const (
	VerifyEmailRoute string = "/api/v1/auth/verify/email/"
)

// Config is a config
type Config struct {
	HttpAddr  string
	SecretKey string

	AllVerifiedByDefault bool

	MongoUrl string
	RedisUrl string

	HealthcheckMessage string

	ClientOrigin string

	AccessPrivateKey   string
	AccessPublicKey    string
	AccessTokenExpired time.Duration
	AccessTokenMaxAge  int

	RefreshPrivateKey   string
	RefreshPublicKey    string
	RefreshTokenExpired time.Duration
	RefreshTokenMaxAge  int

	EmailFrom string
	SMTPPass  string
	SMTPUser  string
	SMTPHost  string
	SMTPPort  string
}

// ReadEnv Read reads config from environment.
func ReadEnv() Config {
	// loads values from .env into the system
	if err := godotenv.Load("dev.env"); err != nil {
		log.Print("No .env file found")
	}
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}

	var config Config
	if httpAddr, exists := os.LookupEnv("HTTP_ADDR"); exists {
		config.HttpAddr = httpAddr
	}
	if allVerifiedByDefault, exists := os.LookupEnv("ALL_VERIFIED_BY_DEFAULT"); exists {
		value, err := strconv.ParseBool(allVerifiedByDefault)
		if err == nil {
			config.AllVerifiedByDefault = value
		}
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
	if accessPrivateKey, exists := os.LookupEnv("ACCESS_TOKEN_PRIVATE_KEY"); exists {
		config.AccessPrivateKey = accessPrivateKey
	}
	if accessPublicKey, exists := os.LookupEnv("ACCESS_TOKEN_PUBLIC_KEY"); exists {
		config.AccessPublicKey = accessPublicKey
	}
	if accessTokenExpired, exists := os.LookupEnv("ACCESS_TOKEN_EXPIRED_IN"); exists {
		value, err := time.ParseDuration(accessTokenExpired)
		if err == nil {
			config.AccessTokenExpired = value
		}
	}
	if accessTokenMaxage, exists := os.LookupEnv("ACCESS_TOKEN_MAXAGE"); exists {
		value, _ := strconv.Atoi(accessTokenMaxage)
		config.AccessTokenMaxAge = value
	}

	if refreshPrivateKey, exists := os.LookupEnv("REFRESH_TOKEN_PRIVATE_KEY"); exists {
		config.RefreshPrivateKey = refreshPrivateKey
	}
	if refreshPublicKey, exists := os.LookupEnv("REFRESH_TOKEN_PUBLIC_KEY"); exists {
		config.RefreshPublicKey = refreshPublicKey
	}
	if refreshTokenExpired, exists := os.LookupEnv("REFRESH_TOKEN_EXPIRED_IN"); exists {
		value, err := time.ParseDuration(refreshTokenExpired)
		if err == nil {
			config.RefreshTokenExpired = value
		}
	}
	if refreshTokenMaxage, exists := os.LookupEnv("REFRESH_TOKEN_MAXAGE"); exists {
		value, _ := strconv.Atoi(refreshTokenMaxage)
		config.RefreshTokenMaxAge = value
	}

	if EmailFrom, exists := os.LookupEnv("EMAIL_FROM"); exists {
		config.EmailFrom = EmailFrom
	}
	if SMTPHost, exists := os.LookupEnv("SMTP_HOST"); exists {
		config.SMTPHost = SMTPHost
	}
	if SMTPPort, exists := os.LookupEnv("SMTP_PORT"); exists {
		config.SMTPPort = SMTPPort
	}
	if SMTPUser, exists := os.LookupEnv("SMTP_USER"); exists {
		config.SMTPUser = SMTPUser
	}
	if SMTPPass, exists := os.LookupEnv("SMTP_PASS"); exists {
		config.SMTPPass = SMTPPass
	}
	return config
}
