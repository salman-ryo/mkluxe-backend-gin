package main

import (
	"log"
	"os"

	"mkluxe-backend/internal/routes"

	"github.com/joho/godotenv"
)

func main() {
	// Load evn file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file was found, hence loading failed.")
	}

	// Port from env and default to 8080 if not set
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Initialize the router
	router := routes.SetupRouter()

	// Start the server
	log.Printf("Starting server on port %s...", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
