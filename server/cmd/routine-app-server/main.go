package main

import (
	"log"
	"routine-app-server/internal/app"

	"github.com/joho/godotenv"
)

// @title Routine App API
// @version 1.0
// @description This is the backend API for the Routine App.
// @BasePath /api/v1
func main() {
	// Read env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Start Application
	app.Run()
}
