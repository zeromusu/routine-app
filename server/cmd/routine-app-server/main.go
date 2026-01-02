package main

import (
	"log"
	"routine-app-server/internal/config"
	"routine-app-server/internal/db"

	"github.com/joho/godotenv"
)

func main() {
	// Read env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Load Config
	config := config.LoadConfig()

	// Create PostgreSQL connection
	_, err := db.InitDB(config.DB)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
}
