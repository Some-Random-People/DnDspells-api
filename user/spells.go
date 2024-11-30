package user

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/some-random-people/dndspells-api/auth"
)

type UserSpell struct {
	Name        string  `json:"name"`
	Level       *int    `json:"level"`
	School      *int    `json:"school"`
	IsRitual    *int    `json:"isRitual"`
	CastingTime *string `json:"castingTime"`
	SpellRange  *string `json:"spellRange"`
	Components  *string `json:"components"`
	Duration    *string `json:"duration"`
	Description *string `json:"description"`
	Upcast      *string `json:"upcast"`
	IsPublic    int     `json:"isPublic"`
}

func CreateUserSpellsEndpoints(router *mux.Router, db *sql.DB) {
	router.HandleFunc("/api/user/spell/create", func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		verified, claims := auth.VerifyToken(token)
		if !verified {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "Unauthorized")
			return
		}
		var newSpell UserSpell
		err := json.NewDecoder(r.Body).Decode(&newSpell)
		if err != nil {
			log.Println("Bad body")
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "Bad Request")
			return
		}
		if newSpell.Name == "" || newSpell.IsPublic == 0 || newSpell.IsRitual == nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "Bad Request")
			return
		}
		spellInsert, err := db.Prepare("INSERT INTO user_spells(`name`, `level`, `school`, `is_ritual`, `casting_time`, `range`, `components`, `duration`, `description`, `upcast`, `is_public`, `user_id`) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
		if err != nil {
			log.Fatal(err)
		}
		res, err := spellInsert.Exec(newSpell.Name, newSpell.Level, newSpell.School, newSpell.IsRitual, newSpell.CastingTime, newSpell.SpellRange, newSpell.Components, newSpell.Duration, newSpell.Description, newSpell.Upcast, newSpell.IsPublic, claims["identifier"])
		if err != nil {
			log.Fatal(err)
		}
		id, err := res.LastInsertId()
		if err != nil {
			log.Fatal(err)
		}
		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"id":%s}`, strconv.FormatInt(id, 10))
	}).Methods("POST")
}
