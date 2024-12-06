package search

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/some-random-people/dndspells-api/dataStructs"
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

		var queryData dataStructs.QueryData
		castingTimesPre := []string{"1_action", "1_bonus_action", "1_hour", "1_minute", "1_reaction", "10_minutes", "12_hours", "24_hours", "8_hours"}
		rangeTypePre := []string{"mile", "feet", "Self", "Sight", "Special", "Touch", "Unlimited"}
		componentsPre := []string{"V", "S", "M", "VS", "VM", "VSM", "SM"}
		durationPre := []string{"1_day", "1_hour", "1_minute", "1_round", "10_days", "10_minutes", "24_hours", "30_days", "7_days", "8_hours", "Concentration", "Instantaneous", "Instantanous", "Special", "Until_dispelled", "Up_to_1_minute", "Up_to_1_hour", "Up_to_8_hours"}

		if r.URL.Query().Get("id") != "" {
			temp, err := strconv.Atoi(r.URL.Query().Get("id"))
			if err != nil {
				log.Printf("Something went wrong with id conversion")
			}
			queryData.Id = temp
		}

		if r.URL.Query().Get("source") != "" {
			temp, err := strconv.Atoi(r.URL.Query().Get("source"))
			if err != nil {
				log.Printf("Something went wrong with source conversion")
			}
			queryData.Source = temp
		}

		if r.URL.Query().Get("level") != "" {
			temp, err := strconv.Atoi(r.URL.Query().Get("level"))
			if err != nil {
				log.Printf("Something went wrong with level conversion")
			}
			queryData.Level = temp
		}

		if r.URL.Query().Get("school") != "" {
			temp, err := strconv.Atoi(r.URL.Query().Get("school"))
			if err != nil {
				log.Printf("Something went wrong with school conversion")
			}
			queryData.School = temp
		}

		if r.URL.Query().Get("isRitual") != "" {
			queryData.IsRitual = true
		} else {
			queryData.IsRitual = false
		}

		if r.URL.Query().Get("castingTime") != "" {
			temp := r.URL.Query().Get("castingTime")
			for _, v := range castingTimesPre {
				if v == temp {
					queryData.CastingTime = v
				}
			}
		}

		if r.URL.Query().Get("rangeValueStart") != "" {
			temp, err := strconv.Atoi(r.URL.Query().Get("rangeValueStart"))
			if err != nil {
				log.Printf("Something went wrong with rangeValueStart conversion")
			}
			queryData.RangeValueStart = temp
		}

		if r.URL.Query().Get("rangeValueStop") != "" {
			temp, err := strconv.Atoi(r.URL.Query().Get("rangeValueStop"))
			if err != nil {
				log.Printf("Something went wrong with rangeValueStop conversion")
			}
			queryData.RangeValueStop = temp
		}

		if r.URL.Query().Get("rangeType") != "" {
			temp := r.URL.Query().Get("rangeType")
			for _, v := range rangeTypePre {
				if v == temp {
					queryData.RangeType = v
				}
			}
		}

		if r.URL.Query().Get("components") != "" {
			temp := r.URL.Query().Get("components")
			for _, v := range componentsPre {
				if v == temp {
					queryData.Components = v
				}
			}
		}

		if r.URL.Query().Get("duration") != "" {
			temp := r.URL.Query().Get("duration")
			for _, v := range durationPre {
				if v == temp {
					queryData.Duration = v
				}
			}
		}

		if r.URL.Query().Get("upcast") != "" {
			queryData.Upcast = true
		} else {
			queryData.Upcast = false
		}

		result, err := DataQuery(db, queryData)
		if err != nil {
			log.Println(err)
		}

		fmt.Fprint(w, string(result))
	})
}
