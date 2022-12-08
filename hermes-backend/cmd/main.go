package main

import (
	"fmt"

	"woojiahao.com/hermes/internal/database"
)

func main() {
	databaseConfiguration := database.LoadConfiguration()
	db := database.Initialize(databaseConfiguration)
	mary, err := db.GetUser("maryanne")
	if err != nil {
		panic(err)
	}

	thread, err := db.CreateThread(mary.Id, "Hello world", "Lorem ipsum dolor", []database.Tag{
		database.NewTag("productivity", "#111111"),
		database.NewTag("fitness", "#188563"),
		database.NewTag("something", "#188563"),
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(thread.Tags)

	thread, err = db.EditThread(mary.Id, thread.Id, "Hello world", "Loream ipsum dolor", false, false, []database.Tag{
		database.NewTag("productivity", "#1111"),
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(thread.Tags)
}
