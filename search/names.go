package search

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func SearchLists(router *mux.Router, db *sql.DB) {

	router.HandleFunc("/api/search/ids", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprint(w, string(SpellList(db)))
	})
}
