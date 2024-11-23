package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/some-random-people/dndspells-api/auth"
	"github.com/some-random-people/dndspells-api/database"
	"github.com/some-random-people/dndspells-api/search"

	"github.com/gorilla/mux"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error while loading .env file")
	}
	log.Println("Running server")
	db := database.ConnectToDatabase()

	router := mux.NewRouter()

	createEndpoint(router)
	auth.DiscordConfig(router, db)
	search.SearchLists(router, db)
	http.ListenAndServe(":80", router)
}
