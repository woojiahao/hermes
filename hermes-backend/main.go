package main

import (
	"log"

	"github.com/joho/godotenv"
	"woojiahao.com/hermes/internal/database"
	"woojiahao.com/hermes/internal/server"
)

// Loads any configurations from .env
// TODO: Support development vs production .env
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

	// Start the server
	serverConfiguration := server.LoadConfiguration()
	server.Start(serverConfiguration, db)
}
