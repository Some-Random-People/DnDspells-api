package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func createEndpoint(router *mux.Router) {
	router.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello world!")
	})
}
