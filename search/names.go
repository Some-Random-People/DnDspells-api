package search

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func SearchLists(router *mux.Router, db *sql.DB) {

	router.HandleFunc("/api/search/spellName", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprint(w, string(SpellList(db)))
	})

	router.HandleFunc("/api/search/schoolName", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprint(w, string(SchoolList(db)))
	})

	router.HandleFunc("/api/search/bookName", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprint(w, string(BookList(db)))
	})

	router.HandleFunc("/api/search/misc", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprint(w, string(Misc(db)))
	})
}
