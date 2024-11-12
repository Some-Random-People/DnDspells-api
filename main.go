package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func main() {
	fmt.Printf("Running server")

	router := mux.NewRouter()

	createEndpoint(router)

	http.ListenAndServe(":80", router)
}
