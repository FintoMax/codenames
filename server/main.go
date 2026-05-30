package main

import (
	"codenames/server/config"
	"codenames/server/database"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatal(err)
	}
	db, err := database.NewDB(cfg)
	defer db.Close()
	if err != nil {
		log.Fatal(err)
	}
}
