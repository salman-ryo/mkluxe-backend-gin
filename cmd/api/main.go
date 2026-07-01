// Main server entry point to connect to Mongo, run the seeder, and boot the newly wired app
package main

import (
	"context"
	"log"
	"os"

	"mkluxe-backend/internal/app"
	"mkluxe-backend/internal/db"
	"mkluxe-backend/internal/seed"

	"github.com/joho/godotenv"
)

func main() {
	// 1. Load config
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "mkluxe"
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// 2. Connect Database
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

	// 3. Ensure Indexes
	if err := db.EnsureIndexes(database); err != nil {
		log.Printf("Warning: Failed to ensure indexes: %v", err)
	}

	// 4. Seed Admin
	if err := seed.SeedSuperAdmin(database); err != nil {
		log.Fatalf("Failed to seed admin: %v", err)
	}

	// 5. Wire App & Router
	router := app.BuildApp(database)

	// 6. Start Server
	log.Printf("Starting MKLuxe server on port %s...", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
}
