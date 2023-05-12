package main

import (
	"log"

	"github.com/alkurbatov/goph-keeper/internal/keeper/app"
	"github.com/alkurbatov/goph-keeper/internal/keeper/config"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	if err := app.Run(cfg); err != nil {
		log.Fatal(err)
	}
}
