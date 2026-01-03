package main

import (
	"log"
	"routine-app-server/internal/app"

	"github.com/joho/godotenv"
)

func main() {
	// Read env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Start Application
	app.Run()
}
