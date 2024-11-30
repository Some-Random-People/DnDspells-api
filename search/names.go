package search

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func SearchLists(router *mux.Router, db *sql.DB) {

	router.HandleFunc("/api/spell/search/spellName", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprint(w, string(SpellList(db)))
	})

	router.HandleFunc("/api/spell/search/schoolName", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprint(w, string(SchoolList(db)))
	})

	router.HandleFunc("/api/spell/search/bookName", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprint(w, string(BookList(db)))
	})

	router.HandleFunc("/api/spell/search/misc", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprint(w, string(Misc(db)))
	})

	router.HandleFunc("/api/spell/search/spell", func(w http.ResponseWriter, r *http.Request) {
		type qd struct { // qd - Query data
			id              int
			source          int
			level           int
			school          int
			isRitual        bool
			castingTime     string // Predefined
			rangeValueStart int
			rangeValueStop  int
			rangeType       string // Predefined OR feet OR mile
			components      string // For reconsideration about selection REGEX FUCK ME
			duration        string // Predefined
			upcast          bool
		}
		var queryData qd
		castingTimesPre := []string{"1_action", "1_bonus_action", "1_hour", "1_minute", "1_reaction", "10_minutes", "12_hours", "24_hours", "8_hours"}
		rangeTypePre := []string{"mile", "feet", "Self", "Sight", "Special", "Touch", "Unlimited"}
		componentsPre := []string{"V", "S", "M", "VS", "VM", "VSM", "SM"}
		durationPre := []string{"1_day", "1_hour", "1_minute", "1_round", "10_days", "10_minutes", "24_hours", "30_days", "7_days", "8_hours", "Concentration", "Instantaneous", "Instantanous", "Special", "Until_dispelled", "Up_to_1_minute", "Up_to_1_hour", "Up_to_8_hours"}

		if r.URL.Query().Get("id") != "" {
			temp, err := strconv.Atoi(r.URL.Query().Get("id"))
			if err != nil {
				log.Printf("Something went wrong with id conversion")
			}
			queryData.id = temp
		}

		if r.URL.Query().Get("source") != "" {
			temp, err := strconv.Atoi(r.URL.Query().Get("source"))
			if err != nil {
				log.Printf("Something went wrong with source conversion")
			}
			queryData.source = temp
		}

		if r.URL.Query().Get("level") != "" {
			temp, err := strconv.Atoi(r.URL.Query().Get("level"))
			if err != nil {
				log.Printf("Something went wrong with level conversion")
			}
			queryData.level = temp
		}

		if r.URL.Query().Get("school") != "" {
			temp, err := strconv.Atoi(r.URL.Query().Get("school"))
			if err != nil {
				log.Printf("Something went wrong with school conversion")
			}
			queryData.school = temp
		}

		if r.URL.Query().Get("isRitual") != "" {
			queryData.isRitual = true
		} else {
			queryData.isRitual = false
		}

		if r.URL.Query().Get("castingTime") != "" {
			temp := r.URL.Query().Get("castingTime")
			for _, v := range castingTimesPre {
				if v == temp {
					queryData.castingTime = v
				}
			}
		}

		if r.URL.Query().Get("rangeValueStart") != "" {
			temp, err := strconv.Atoi(r.URL.Query().Get("rangeValueStart"))
			if err != nil {
				log.Printf("Something went wrong with rangeValueStart conversion")
			}
			queryData.rangeValueStart = temp
		}

		if r.URL.Query().Get("rangeValueStop") != "" {
			temp, err := strconv.Atoi(r.URL.Query().Get("rangeValueStop"))
			if err != nil {
				log.Printf("Something went wrong with rangeValueStop conversion")
			}
			queryData.rangeValueStop = temp
		}

		if r.URL.Query().Get("rangeType") != "" {
			temp := r.URL.Query().Get("rangeType")
			for _, v := range rangeTypePre {
				if v == temp {
					queryData.rangeType = v
				}
			}
		}

		if r.URL.Query().Get("components") != "" {
			temp := r.URL.Query().Get("components")
			for _, v := range componentsPre {
				if v == temp {
					queryData.components = v
				}
			}
		}

		if r.URL.Query().Get("duration") != "" {
			temp := r.URL.Query().Get("duration")
			for _, v := range durationPre {
				if v == temp {
					queryData.duration = v
				}
			}
		}

		if r.URL.Query().Get("upcast") != "" {
			queryData.upcast = true
		} else {
			queryData.upcast = false
		}

		log.Println(queryData)
	})
}
