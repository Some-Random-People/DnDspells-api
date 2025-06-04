package user

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/some-random-people/dndspells-api/auth"
	"github.com/some-random-people/dndspells-api/dataStructs"
	"github.com/some-random-people/dndspells-api/utils"
)

func CreateUserProfileEndpoints(router *mux.Router, db *sql.DB) {
	// Retrive user profile
	router.HandleFunc("/api/user/profile/{id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		var profile dataStructs.UserProfile
		rows, err := db.Query("SELECT `id`, `nickname` FROM users WHERE id = ?", id)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer rows.Close()
		var exist = false
		for rows.Next() {
			if err := rows.Scan(&profile.Id, &profile.Name); err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			exist = true
		}
		if exist == false {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "User not found")
			return
		}
		jsonData, err := json.Marshal(profile)
		if err != nil {
			log.Printf("Can't convert profile to json %s\n", err)
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, string(jsonData))
	})

	// Retrive user spells
	router.HandleFunc("/api/user/profile/{id}/spells", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		var spells []dataStructs.UserSpell
		rows, err := db.Query("SELECT `id`, `name`, `level`, `school`, `is_ritual`, `casting_time`, `range`, `components`, `duration`, `description`, `upcast`, `user_id`, `is_public` FROM user_spells WHERE user_id = ? AND is_public = 1", id)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer rows.Close()
		var exist = false
		for rows.Next() {
			var spell dataStructs.UserSpell
			if err := rows.Scan(&spell.Id, &spell.Name, &spell.Level, &spell.School, &spell.IsRitual, &spell.CastingTime, &spell.SpellRange, &spell.Components, &spell.Duration, &spell.Description, &spell.Upcast, &spell.User_id, &spell.IsPublic); err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			spells = append(spells, spell)
			exist = true
		}
		if exist == false {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "User not found")
			return
		}

		jsonData, err := json.Marshal(spells)
		if err != nil {
			log.Printf("Can't convert user spells to json %s\n", err)
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, string(jsonData))
	})

	router.HandleFunc("/api/user/profile/", func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		verified, claims := auth.VerifyToken(token)
		if !verified {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "Unauthorized")
			return
		}

		var modifiedUser dataStructs.UserProfile
		contentType := r.Header.Get("Content-Type")
		if strings.Contains(contentType, "multipart/form-data") || strings.Contains(contentType, "application/x-www-form-urlencoded") {
			err := utils.ParseForm(&modifiedUser, r)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprint(w, "Bad Request")
				return
			}
		}
		rows, err := db.Query("SELECT `nickname` FROM users WHERE nickname = ?", modifiedUser.Name)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer rows.Close()
		var exist = false
		for rows.Next() {
			exist = true
		}
		if exist == true {
			w.WriteHeader(http.StatusConflict)
			fmt.Fprint(w, "User with that name already exists")
			return
		}
		userUpdate, err := db.Prepare("UPDATE `users` SET `nickname` = ? WHERE `id` = ?;")
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		_, err = userUpdate.Exec(modifiedUser.Name, claims["identifier"])
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}).Methods("PATCH")
}
