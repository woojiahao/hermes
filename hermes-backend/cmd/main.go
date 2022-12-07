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

	db.CreateUser("johndoe", "johndoe@gmail.com", "1234")
	db.CreateUser("maryanne", "maryanne@gmail.com", "helloworld")
	user, _ := db.GetUser("johndoe")
	fmt.Println(user.Email)

	users, _ := db.GetUsers()
	for _, user := range users {
		fmt.Printf("username: %s, email: %s", user.Username, user.Email)
	}
}
