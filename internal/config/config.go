package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Server struct {
		Port string
	}

	Redis struct {
		Addr     string
		Password string
		DB       int
	}

	Email struct {
		APIKey     string
		FromEmail  string
		Endpoint   string
		CodeLength int
		CodeExpiry time.Duration
	}
}

func LoadConfig() (*Config, error) {
	// .env файл
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	cfg := &Config{}

	// Server config
	cfg.Server.Port = os.Getenv("SERVER_PORT")
	if cfg.Server.Port == "" {
		cfg.Server.Port = "8080"
	}

	// Redis config
	cfg.Redis.Addr = os.Getenv("REDIS_ADDR")
	if cfg.Redis.Addr == "" {
		cfg.Redis.Addr = "localhost:6379"
	}

	cfg.Redis.Password = os.Getenv("REDIS_PASSWORD")

	db, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		db = 0
	}
	cfg.Redis.DB = db

	// Email config
	cfg.Email.APIKey = os.Getenv("ELASTIC_EMAIL_API_KEY")
	cfg.Email.FromEmail = os.Getenv("ELASTIC_EMAIL_FROM")
	cfg.Email.Endpoint = os.Getenv("ELASTIC_EMAIL_ENDPOINT")

	codeLength, err := strconv.Atoi(os.Getenv("VERIFICATION_CODE_LENGTH"))
	if err != nil {
		codeLength = 5
	}
	cfg.Email.CodeLength = codeLength

	expiryMinutes, err := strconv.Atoi(os.Getenv("VERIFICATION_CODE_EXPIRY_MINUTES"))
	if err != nil {
		expiryMinutes = 10
	}
	cfg.Email.CodeExpiry = time.Duration(expiryMinutes) * time.Minute

	return cfg, nil
}
