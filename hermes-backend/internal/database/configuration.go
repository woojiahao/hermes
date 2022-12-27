package database

import (
	"log"
	"os"
	"strconv"
)

type Configuration struct {
	Username string
	Password string
	Host     string
	Name     string
	Port     int
}

func LoadConfiguration() *Configuration {
	username := os.Getenv("DATABASE_USERNAME")
	password := os.Getenv("DATABASE_PASSWORD")
	host := os.Getenv("DATABASE_HOST")
	name := os.Getenv("DATABASE_NAME")
	port, err := strconv.Atoi(os.Getenv("DATABASE_PORT"))
	if err != nil {
		log.Fatal("Invalid DATABASE_PORT in .env")
	}

	return &Configuration{username, password, host, name, port}
}
