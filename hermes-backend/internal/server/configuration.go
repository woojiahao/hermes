package server

import "os"

type Configuration struct {
	JWTKey string
}

func LoadConfiguration() *Configuration {
	return &Configuration{os.Getenv("JWT_KEY")}
}
