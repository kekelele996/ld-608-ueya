package config

import (
	"fmt"
	"os"
)

// Config carries runtime settings sourced from environment variables so that
// .env.example / docker-compose.yml / config / request wrappers stay coupled.
type Config struct {
	Port        string
	DBHost      string
	DBPort      string
	DBName      string
	DBUser      string
	DBPassword  string
	SQLitePath  string
	JWTSecret   string
	SeedOnStart bool
}

func getenv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// Load reads the process environment and never hardcodes connection details.
func Load() *Config {
	return &Config{
		Port:        getenv("PORT", "3000"),
		DBHost:      getenv("DB_HOST", "127.0.0.1"),
		DBPort:      getenv("DB_PORT", "3306"),
		DBName:      getenv("DB_NAME", "app_db"),
		DBUser:      getenv("DB_USER", "app_user"),
		DBPassword:  getenv("DB_PASSWORD", "app_password"),
		SQLitePath:  getenv("SQLITE_PATH", ""),
		JWTSecret:   getenv("JWT_SECRET", "local-dev-secret"),
		SeedOnStart: getenv("SEED_ON_START", "true") == "true",
	}
}

// DSN builds the MySQL DSN consumed by GORM.
func (c *Config) DSN() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName,
	)
}
