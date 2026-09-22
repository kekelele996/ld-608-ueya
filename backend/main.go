package main

import (
	"log"

	"groundTurn/src/config"
	"groundTurn/src/database"
	"groundTurn/src/routes"
	"groundTurn/src/services"
)

func main() {
	cfg := config.Load()

	db, err := database.Open(cfg)
	if err != nil {
		log.Fatalf("connect database: %v", err)
	}
	if cfg.SeedOnStart {
		if err := database.Seed(db); err != nil {
			log.Fatalf("seed database: %v", err)
		}
	}

	svc := services.NewServices(db, cfg)
	if err := svc.Auth.EnsureSeedUsers(); err != nil {
		log.Fatalf("ensure users: %v", err)
	}

	addr := ":" + cfg.Port
	log.Printf("ground-turn backend listening on %s", addr)
	if err := routes.Start(addr, cfg, svc); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}
