package search

import (
	"database/sql"
	"encoding/json"
	"log"
	"strconv"

	"github.com/some-random-people/dndspells-api/dataStructs"
)

func DataQuery(db *sql.DB, queryData dataStructs.QueryData) (string, error) {
	type qr struct { // qr - Query result
		Id          int
		Name        string
		Source      string
		Level       int
		School      string
		IsRitual    bool
		CastingTime string
		SpellRange  string
		Components  string
		Duration    string
		Description string
		Upcast      *string
	}

	var (
		queryResult qr
		//queryDataString dataStructs.QueryDataStrings
	)

	stmt, err := db.Prepare("SELECT basic_spells.id, basic_spells.`name`, books.`name`, basic_spells.`level`, schools.`name`, basic_spells.is_ritual, basic_spells.casting_time, basic_spells.`range`, basic_spells.components, basic_spells.duration, basic_spells.`description`, basic_spells.upcast FROM basic_spells JOIN books ON basic_spells.`source` LIKE books.id JOIN schools ON basic_spells.school LIKE schools.id WHERE basic_spells.id LIKE ? AND basic_spells.`source` LIKE ? AND basic_spells.level LIKE ? AND basic_spells.school LIKE ? AND basic_spells.is_ritual LIKE ? AND basic_spells.casting_time LIKE ? AND basic_spells.range LIKE ? AND basic_spells.components LIKE ? AND basic_spells.duration LIKE ? AND basic_spells.description LIKE ? AND COALESCE(basic_spells.upcast, '') LIKE ?")

	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	log.Println(queryData.Id)
	log.Println(strconv.Itoa(queryData.Id))

	rows, err := stmt.Query("1", "%", "%", "%", "%", "%", "%", "%", "%", "%", "%")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var results []qr
	for rows.Next() {
		err := rows.Scan(&queryResult.Id, &queryResult.Name, &queryResult.Source, &queryResult.Level, &queryResult.School, &queryResult.IsRitual, &queryResult.CastingTime, &queryResult.SpellRange, &queryResult.Components, &queryResult.Duration, &queryResult.Description, &queryResult.Upcast)
		if err != nil {
			log.Fatal(err)
		}
		results = append(results, queryResult)
	}
	log.Println(results)

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	resultsJSON, err := json.Marshal(results)
	if err != nil {
		log.Println(err)
		return "", err
	}

	return string(resultsJSON), nil
}
