package config

import (
	"os"

	_ "github.com/joho/godotenv/autoload" // load .env file automatically
)

type Config struct {
	Port       string
	SessionKey string
}

func Load() Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "1323"
	}
	secret := os.Getenv("SESSION_KEY")
	if secret == "" {
		secret = "gaspoll-local-secret"
	}
	return Config{
		Port:       port,
		SessionKey: secret,
	}
}
