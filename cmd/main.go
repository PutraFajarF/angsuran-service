package main

import (
	"angsuran-service/config"
	"angsuran-service/internal/app"
	"log"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	log.Println("start running app")

	app.Run(cfg)
}
