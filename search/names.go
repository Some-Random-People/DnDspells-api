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

		w.Header().Add("Content-Type", "application/json")
		var queryData dataStructs.QueryData
		castingTimesPre := []string{"1_action", "1_bonus_action", "1_hour", "1_minute", "1_reaction", "10_minutes", "12_hours", "24_hours", "8_hours"}
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
			queryData.Level = temp + 1
		}

		if r.URL.Query().Get("school") != "" {
			temp, err := strconv.Atoi(r.URL.Query().Get("school"))
			if err != nil {
				log.Printf("Something went wrong with school conversion")
			}
			queryData.School = temp
		}

		if r.URL.Query().Get("isRitual") == "t" {
			queryData.IsRitual = "1"
		} else if r.URL.Query().Get("isRitual") == "f" {
			queryData.IsRitual = "0"
		} else {
			queryData.IsRitual = "%"
		}

		queryData.CastingTime = "%"
		if r.URL.Query().Get("castingTime") != "" {
			temp := r.URL.Query().Get("castingTime")
			for _, v := range castingTimesPre {
				if v == temp {
					queryData.CastingTime = v
				}
			}
		}

		// Different values for feet and miles and specials
		// rvsf = range value start feet
		// rvef = range value end feet
		// rvsm = range value start miles
		// rvem = range value end miles
		// rs = range specials
		//
		// Every value is int
		// min 0 max 5000
		// For rs is like binary
		// 1 = Self
		// 2 = Sight
		// 4 = Special
		// 8 = Touch
		// 16 = Unlimited
		// So 31 is like % and is a max value
		// 0 means none of the specials
		//

		queryData.RVSF = 0
		if r.URL.Query().Get("rvsf") != "" {
			temp, err := strconv.Atoi(r.URL.Query().Get("rvsf"))
			if err != nil {
				log.Printf("Something went wrong with rvsf conversion")
			}
			queryData.RVSF = temp
		}

		queryData.RVEF = 5000
		if r.URL.Query().Get("rvef") != "" {
			temp, err := strconv.Atoi(r.URL.Query().Get("rvef"))
			if err != nil {
				log.Printf("Something went wrong with rvef conversion")
			}
			queryData.RVEF = temp
		}

		queryData.RVSM = 0
		if r.URL.Query().Get("rvsm") != "" {
			temp, err := strconv.Atoi(r.URL.Query().Get("rvsm"))
			if err != nil {
				log.Printf("Something went wrong with rvsm conversion")
			}
			queryData.RVSM = temp
		}

		queryData.RVEM = 5000
		if r.URL.Query().Get("rvem") != "" {
			temp, err := strconv.Atoi(r.URL.Query().Get("rvem"))
			if err != nil {
				log.Printf("Something went wrong with rvem conversion")
			}
			queryData.RVEM = temp
		}

		if r.URL.Query().Get("rs") != "" {
			queryData.RS = r.URL.Query().Get("rs")
		} else {
			queryData.RS = "11111"
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
