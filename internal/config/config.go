package config

import (
	"os"
	"strconv"
	"time"
)

// Config bundles runtime configuration loaded from environment variables.
type Config struct {
	AppName        string
	AppEnv         string
	HTTPPort       string
	DatabaseURL    string
	SessionSecret  string
	JWTSigningKey  string
	UptraceDSN     string
	PaymentBaseURL string
	RewardTarget   int
	RewardValue    int
	TokenTTL       time.Duration
}

// Load pulls configuration from environment variables with sensible defaults
// so the app can be bootstrapped quickly in local/dev environments.
func Load() Config {
	return Config{
		AppName:        getEnv("APP_NAME", "Gaspoll"),
		AppEnv:         getEnv("APP_ENV", "local"),
		HTTPPort:       getEnv("HTTP_PORT", "1323"),
		DatabaseURL:    getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/gaspoll?sslmode=disable"),
		SessionSecret:  getEnv("SESSION_SECRET", "dev-secret-please-change"),
		JWTSigningKey:  getEnv("JWT_SIGNING_KEY", "dev-jwt-signing-key"),
		UptraceDSN:     getEnv("UPTRACE_DSN", ""),
		PaymentBaseURL: getEnv("PAYMENT_BASE_URL", "https://sandbox-payments.invalid"),
		RewardTarget:   getEnvInt("REWARD_TARGET", 500000), // nominal akumulasi dalam rupiah
		RewardValue:    getEnvInt("REWARD_VALUE", 75000),
		TokenTTL:       getEnvDuration("TOKEN_TTL", 24*time.Hour),
	}
}

func getEnv(key, def string) string {
	val := os.Getenv(key)
	if val == "" {
		return def
	}
	return val
}

func getEnvInt(key string, def int) int {
	val := os.Getenv(key)
	if val == "" {
		return def
	}
	intVal, err := strconv.Atoi(val)
	if err != nil {
		return def
	}
	return intVal
}

func getEnvDuration(key string, def time.Duration) time.Duration {
	val := os.Getenv(key)
	if val == "" {
		return def
	}
	d, err := time.ParseDuration(val)
	if err != nil {
		return def
	}
	return d
}
