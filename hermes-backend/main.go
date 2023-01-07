package main

import (
	"log"

	"github.com/joho/godotenv"
	"woojiahao.com/hermes/internal/database"
	"woojiahao.com/hermes/internal/server"
)

// Loads any configurations from .env
func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	loadEnv()

	// Setup the database connection
	databaseConfiguration := database.LoadConfiguration()
	db := database.Initialize(databaseConfiguration)
	db.CreateTablesIfNotExists()

	// Setup the original admin user for testing purposes
	if _, err := db.GetUser("admin"); err != nil {
		_, err = db.CreateUser("admin", "root", database.ADMIN)
		if err != nil {
			log.Fatalf("Failed to setup admin user %s", err)
		}
	}

	// Start the server
	serverConfiguration := server.LoadConfiguration()
	server.Start(serverConfiguration, db)
}
