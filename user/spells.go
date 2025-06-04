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

		// Data integrity checks

		if *newSpell.IsPublic != 0 && *newSpell.IsPublic != 1 {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "Bad Request")
			return
		}

		if *newSpell.IsRitual != 0 && *newSpell.IsRitual != 1 {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "Bad Request")
			return
		}

		if newSpell.School != nil {
			if *newSpell.School < 1 || *newSpell.School > 8 {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprint(w, "Bad Request")
				return
			}
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

	// Updating spell
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
		var user_id int
		rows, err := db.Query("SELECT `user_id` FROM user_spells WHERE id = ?", &newSpell.Id)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()
		var exist = false
		for rows.Next() {
			if err := rows.Scan(&user_id); err != nil {
				log.Fatal(err)
			}
			exist = true
		}
		if exist == false {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "Spell not found")
			return
		}

		if claims["identifier"] != strconv.Itoa(user_id) {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "Unauthorized")
			return
		}

		tx, err := db.Begin()
		if err != nil {
			log.Fatal(err)
		}
		defer tx.Rollback()

		if newSpell.IsPublic != nil {
			if *newSpell.IsPublic != 0 && *newSpell.IsPublic != 1 {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprint(w, "Bad Request")
				return
			}
			_, err = tx.Exec("UPDATE user_spells SET is_public = ? WHERE id = ?", *newSpell.IsPublic, newSpell.Id)
			if err != nil {
				log.Println(err)
				return
			}
		}

		if newSpell.IsRitual != nil {
			if *newSpell.IsRitual != 0 && *newSpell.IsRitual != 1 {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprint(w, "Bad Request")
				return
			}
			_, err = tx.Exec("UPDATE user_spells SET is_ritual = ? WHERE id = ?", *newSpell.IsRitual, newSpell.Id)
			if err != nil {
				log.Println(err)
				return
			}
		}

		if newSpell.School != nil {
			if *newSpell.School < 1 || *newSpell.School > 8 {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprint(w, "Bad Request")
				return
			}
			_, err = tx.Exec("UPDATE user_spells SET school = ? WHERE id = ?", *newSpell.School, newSpell.Id)
			if err != nil {
				log.Println(err)
				return
			}
		}

		if newSpell.Level != nil {
			_, err = tx.Exec("UPDATE user_spells SET level = ? WHERE id = ?", *newSpell.Level, newSpell.Id)
			if err != nil {
				log.Println(err)
				return
			}
		}

		if newSpell.Name != "" {
			_, err = tx.Exec("UPDATE user_spells SET name = ? WHERE id = ?", newSpell.Name, newSpell.Id)
			if err != nil {
				log.Println(err)
				return
			}
		}

		if newSpell.SpellRange != nil {
			_, err = tx.Exec("UPDATE user_spells SET range = ? WHERE id = ?", *newSpell.SpellRange, newSpell.Id)
			if err != nil {
				log.Println(err)
				return
			}
		}

		if newSpell.CastingTime != nil {
			_, err = tx.Exec("UPDATE user_spells SET casting_time = ? WHERE id = ?", *newSpell.CastingTime, newSpell.Id)
			if err != nil {
				log.Println(err)
				return
			}
		}

		if newSpell.Components != nil {
			_, err = tx.Exec("UPDATE user_spells SET components = ? WHERE id = ?", *newSpell.Components, newSpell.Id)
			if err != nil {
				log.Println(err)
				return
			}
		}

		if newSpell.Duration != nil {
			_, err = tx.Exec("UPDATE user_spells SET duration = ? WHERE id = ?", *newSpell.Duration, newSpell.Id)
			if err != nil {
				log.Println(err)
				return
			}
		}

		if newSpell.Description != nil {
			_, err = tx.Exec("UPDATE user_spells SET description = ? WHERE id = ?", *newSpell.Description, newSpell.Id)
			if err != nil {
				log.Println(err)
				return
			}
		}

		if newSpell.Upcast != nil {
			_, err = tx.Exec("UPDATE user_spells SET upcast = ? WHERE id = ?", *newSpell.Upcast, newSpell.Id)
			if err != nil {
				log.Println(err)
				return
			}
		}

		err = tx.Commit()
		if err != nil {
			log.Println(err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}).Methods("PATCH")

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
			w.WriteHeader(http.StatusNotFound)
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

	// Deleting spell
	router.HandleFunc("/api/user/spell/", func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		verified, claims := auth.VerifyToken(token)
		if !verified {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "Unauthorized")
			return
		}
		var targetSpell dataStructs.UserSpell
		contentType := r.Header.Get("Content-Type")
		if strings.Contains(contentType, "multipart/form-data") || strings.Contains(contentType, "application/x-www-form-urlencoded") {
			err := utils.ParseForm(&targetSpell, r)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprint(w, "Bad Request")
				return
			}
		} else if contentType == "application/json" {
			err := json.NewDecoder(r.Body).Decode(&targetSpell)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprint(w, "Bad Request")
				return
			}
		}
		if targetSpell.Id == 0 {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "Bad Request")
			return
		}

		var selectedSpell dataStructs.UserSpell
		rows, err := db.Query("SELECT `id`, `user_id` FROM user_spells WHERE id = ?", strconv.Itoa(targetSpell.Id))
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()
		var exist = false
		for rows.Next() {
			if err := rows.Scan(&selectedSpell.Id, &selectedSpell.User_id); err != nil {
				log.Fatal(err)
			}
			exist = true
		}
		if exist == false {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "Spell not found")
			return
		}

		if claims["identifier"] != strconv.Itoa(selectedSpell.User_id) {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "Unauthorized")
			return

		}
		spellInsert, err := db.Prepare("DELETE FROM user_spells WHERE id = ?")
		if err != nil {
			log.Fatal(err)
		}
		_, err = spellInsert.Exec(targetSpell.Id)
		if err != nil {
			log.Fatal(err)
		}
		w.WriteHeader(http.StatusNoContent)
	}).Methods("DELETE")
}
