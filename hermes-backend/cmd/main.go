package main

import (
	"fmt"

	"woojiahao.com/hermes/internal/database"
)

func main() {
	databaseConfiguration := database.LoadConfiguration()
	db := database.Initialize(databaseConfiguration)

	db.CreateUser("johndoe", "johndoe@gmail.com", "1234", database.ADMIN)
	_, err := db.CreateUser("maryanne", "maryanne@gmail.com", "helloworld", database.USER)
	if err != nil {
		fmt.Println(err)
	}
	user, _ := db.GetUser("johndoe")
	mary, _ := db.GetUser("maryanne")
	fmt.Println(user.Email)

	users, _ := db.GetUsers()
	for _, user := range users {
		fmt.Printf("username: %s, email: %s", user.Username, user.Email)
	}

	thread, err := db.CreateThread(user.Id, "something", "afoqwefj")
	if err != nil {
		fmt.Println(err)
	}
	_, err = db.DeleteThread(mary.Id, thread.Id)
	if err != nil {
		fmt.Println(mary.Id)
		fmt.Println(thread.Id)
		fmt.Println(err)
	}

	maryThread, _ := db.CreateThread(mary.Id, "foo", "bar")
	db.DeleteThread(user.Id, maryThread.Id)

	db.EditThread(user.Id, maryThread.Id, "wfqoei", "sfqewf", false, false)
}
