package user

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/some-random-people/dndspells-api/auth"
	"github.com/some-random-people/dndspells-api/dataStructs"
	"github.com/some-random-people/dndspells-api/utils"
)

func CreateUserSpellsEndpoints(router *mux.Router, db *sql.DB) {
	// Creating new spell
	router.HandleFunc("/api/user/spell/", func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		verified, claims := auth.VerifyToken(token)
		if !verified {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "Unauthorized")
			return
		}
		var newSpell dataStructs.UserSpell
		contentType := r.Header.Get("Content-Type")
		if strings.Contains(contentType, "multipart/form-data") || strings.Contains(contentType, "application/x-www-form-urlencoded") {
			err := utils.ParseForm(&newSpell, r)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprint(w, "Bad Request")
				return
			}
		} else if contentType == "application/json" {
			err := json.NewDecoder(r.Body).Decode(&newSpell)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprint(w, "Bad Request")
				return
			}
		}
		if newSpell.Name == "" || newSpell.IsPublic == nil || newSpell.IsRitual == nil {
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

	// Getting spell
	router.HandleFunc("/api/user/spell/{id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		var spell dataStructs.UserSpell
		rows, err := db.Query("SELECT `id`, `name`, `level`, `school`, `is_ritual`, `casting_time`, `range`, `components`, `duration`, `description`, `upcast`, `user_id`, `is_public` FROM user_spells WHERE id = ?", id)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()
		var exist = false
		for rows.Next() {
			if err := rows.Scan(&spell.Id, &spell.Name, &spell.Level, &spell.School, &spell.IsRitual, &spell.CastingTime, &spell.SpellRange, &spell.Components, &spell.Duration, &spell.Description, &spell.Upcast, &spell.User_id, &spell.IsPublic); err != nil {
				log.Fatal(err)
			}
			exist = true
		}
		if exist == false {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "Spell not found")
			return
		}
		if *spell.IsPublic == 1 {
			response, err := json.Marshal(spell)
			if err != nil {
				log.Println(err)
				return
			}
			w.Header().Add("Content-Type", "application/json")
			fmt.Fprint(w, string(response))
		} else {
			token := r.Header.Get("Authorization")
			verified, claims := auth.VerifyToken(token)
			if !verified {
				w.WriteHeader(http.StatusUnauthorized)
				fmt.Fprint(w, "Unauthorized")
				return
			}
			if claims["identifier"] == strconv.Itoa(spell.User_id) {
				response, err := json.Marshal(spell)
				if err != nil {
					log.Println(err)
					return
				}
				w.Header().Add("Content-Type", "application/json")
				fmt.Fprint(w, string(response))
				return
			}
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "Unauthorized")
		}
	}).Methods("GET")
}
