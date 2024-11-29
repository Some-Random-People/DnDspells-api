package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/some-random-people/dndspells-api/auth"
)

func createEndpoint(router *mux.Router) {
	router.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		authorized, claims := auth.VerifyToken(token)
		if !authorized {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "Unauthorized")
			return
		}
		fmt.Fprintf(w, "Hello world %s!", claims["identifier"])
	})
}
