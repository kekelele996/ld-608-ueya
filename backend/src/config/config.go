package config

import (
	"os"
	"strconv"
)

// Config is assembled from env only; new settings must be mirrored in
// .env.example, docker-compose.yml and the request/log modules.
type Config struct {
	Port         string
	DBHost       string
	DBPort       string
	DBName       string
	DBUser       string
	DBPassword   string
	JWTSecret    string
	SeedOnStart  bool
	RateLimitQPS int
}

func getenv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func Load() Config {
	return Config{
		Port:         getenv("PORT", "3000"),
		DBHost:       getenv("DB_HOST", "127.0.0.1"),
		DBPort:       getenv("DB_PORT", "3306"),
		DBName:       getenv("DB_NAME", "app_db"),
		DBUser:       getenv("DB_USER", "app_user"),
		DBPassword:   getenv("DB_PASSWORD", "app_password"),
		JWTSecret:    getenv("JWT_SECRET", "local-dev-secret"),
		SeedOnStart:  getenv("SEED_ON_START", "true") == "true",
		RateLimitQPS: atoi(getenv("RATE_LIMIT_QPS", "20")),
	}
}

func atoi(value string) int {
	n, err := strconv.Atoi(value)
	if err != nil {
		return 20
	}
	return n
}

// MySQLDSN builds the GORM DSN for the configured MySQL instance.
func (c Config) MySQLDSN() string {
	return c.DBUser + ":" + c.DBPassword + "@tcp(" + c.DBHost + ":" + c.DBPort + ")/" +
		c.DBName + "?charset=utf8mb4&parseTime=true&loc=Local"
}
