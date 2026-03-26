package main

import (
	"log"
	"net/http"

	"assignment5/internal/app"
	"assignment5/internal/handler"
	"assignment5/internal/repository"
)

func main() {
	db := app.MustConnectDB()
	repo := repository.NewUserRepository(db)
	h := handler.NewUserHandler(repo)

	mux := http.NewServeMux()
	mux.HandleFunc("/users", h.GetUsers)
	mux.HandleFunc("/common-friends", h.GetCommonFriends)

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}