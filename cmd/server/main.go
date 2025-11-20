package main

import (
	"log"

	"gaspoll/internal/config"
	"gaspoll/internal/repository/memory"
	"gaspoll/internal/server"
)

func main() {
	cfg := config.Load()
	repo := memory.NewSeededRepository()

	app := server.New(cfg, repo)
	log.Printf("Gaspoll server listening on :%s", cfg.Port)
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
