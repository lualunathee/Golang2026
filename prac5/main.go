package main

import (
	"log"
	"net/http"
)

func main() {

	db := InitDB()

	repo := NewRepository(db)

	handler := NewHandler(repo)

	http.HandleFunc("/users", handler.GetUsers)
	http.HandleFunc("/common-friends", handler.CommonFriends)

	log.Println("Server running on :8080")

	http.ListenAndServe(":8080", nil)
}
