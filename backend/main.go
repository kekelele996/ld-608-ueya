package main

import (
	"log"
	"time"

	"groundTurn/src/config"
	"groundTurn/src/routes"

	"gorm.io/gorm"
)

// main waits for MySQL (compose healthcheck covers most cases), migrates,
// seeds and starts the HTTP server on :3000 inside the container.
func main() {
	cfg := config.Load()

	gormDB, err := openWithRetry(cfg, 30, 2*time.Second)
	if err != nil {
		log.Fatalf("database unavailable: %v", err)
	}
	if cfg.SeedOnStart {
		if err := config.Seed(gormDB); err != nil {
			log.Printf("seed skipped: %v", err)
		}
	}
	log.Printf("ground-turn backend listening on :%s", cfg.Port)
	if err := routes.Start(":"+cfg.Port, gormDB, cfg); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}

func openWithRetry(cfg config.Config, attempts int, pause time.Duration) (*gorm.DB, error) {
	var lastErr error
	for i := 0; i < attempts; i++ {
		db, err := config.OpenMySQL(cfg)
		if err == nil {
			return db, nil
		}
		lastErr = err
		time.Sleep(pause)
	}
	return nil, lastErr
}
