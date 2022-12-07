package main

import (
	"fmt"
	"log"

	"woojiahao.com/hermes/internal/database"
)

func main() {
	databaseConfiguration := database.LoadConfiguration()
	db := database.Initialize(databaseConfiguration)
	roles, err := db.GetRoles()
	if err != nil {
		log.Fatal(err)
	}

	for _, role := range roles {
		fmt.Println(role.Title)
	}
}
