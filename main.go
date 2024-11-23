package main

import (
	"fmt"
	"github.com/some-random-people/dndspells-api/auth"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	fmt.Printf("Running server")

	router := mux.NewRouter()

	createEndpoint(router)
	auth.DiscordConfig(router)
	http.ListenAndServe(":80", router)
}
