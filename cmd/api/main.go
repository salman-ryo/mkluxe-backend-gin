// Main server entry point to connect to Mongo, run the seeder, and boot the newly wired app
package main

import (
	"context"
	"log"
	"os"

	"mkluxe-backend/internal/app"
	"mkluxe-backend/internal/config"
	"mkluxe-backend/internal/db"
	"mkluxe-backend/internal/seed"

	"github.com/joho/godotenv"
)

func main() {
	// 1. Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// 2. Load centralized configuration
	cfg := config.Load()

	// Handle database specific envs
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "mkluxe"
	}

	// 3. Connect Database
	client, database, err := db.ConnectMongo(mongoURI, dbName)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	defer func() {
		log.Println("Disconnecting from MongoDB...")
		if err := client.Disconnect(context.Background()); err != nil {
			log.Printf("Error disconnecting from MongoDB: %v", err)
		}
	}()

	// 4. Ensure Indexes
	if err := db.EnsureIndexes(database); err != nil {
		log.Printf("Warning: Failed to ensure indexes: %v", err)
	}

	// 5. Seed Admin & Categories
	if err := seed.SeedSuperAdmin(database); err != nil {
		log.Fatalf("Failed to seed admin: %v", err)
	}
	if err := seed.SeedCategories(database); err != nil {
		log.Fatalf("Failed to seed categories: %v", err)
	}

	// 6. Wire App & Router
	// Pass the config struct down into the application builder
	router := app.BuildApp(database, cfg)

	// 7. Start Server using the port from our centralized config
	log.Printf("Starting MKLuxe server on port %s...", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
}
