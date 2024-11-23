package search

import (
	"database/sql"
	"encoding/json"
	"log"
)

func SpellList(db *sql.DB) string {
	var spellList = make(map[string]int)
	var (
		id   int
		name string
	)
	rows, err := db.Query("SELECT id, name FROM table_view")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&id, &name)
		if err != nil {
			log.Fatal(err)
		}
		spellList[name] = id
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	jsonBytes, err := json.Marshal(spellList)
	if err != nil {
		log.Fatal(err)
	}

	return string(jsonBytes)
}
