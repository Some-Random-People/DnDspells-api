package search

import (
	"database/sql"
	"strings"
	//"encoding/csv"
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
		queryResult     qr
		queryDataString dataStructs.QueryDataStrings
	)
	queryDataString.Id = integrityInt(queryData.Id)
	queryDataString.Source = integrityInt(queryData.Source)
	queryDataString.Level = integrityLevel(queryData.Level)
	queryDataString.School = integrityInt(queryData.School)
	queryDataString.IsRitual = queryData.IsRitual
	queryDataString.CastingTime = queryData.CastingTime
	rangeVal := rangeValue(queryData.RVSF, queryData.RVEF, queryData.RVSM, queryData.RVEM)
	queryDataString.RVSF = rangeVal[0]
	queryDataString.RVEF = rangeVal[1]
	queryDataString.RVSM = rangeVal[2]
	queryDataString.RVEM = rangeVal[3]
	queryDataString.Specials = specials(queryData.RS)
	//log.Println(queryDataString.Specials)

	stmt, err := db.Prepare("SELECT basic_spells.id, basic_spells.`name`, books.`name`, basic_spells.`level`, schools.`name`, basic_spells.is_ritual, basic_spells.casting_time, basic_spells.`range`, basic_spells.components, basic_spells.duration, basic_spells.`description`, basic_spells.upcast FROM basic_spells JOIN books ON basic_spells.`source` LIKE books.id JOIN schools ON basic_spells.school LIKE schools.id WHERE basic_spells.id LIKE ? AND basic_spells.`source` LIKE ? AND basic_spells.level LIKE ? AND basic_spells.school LIKE ? AND basic_spells.is_ritual LIKE ? AND basic_spells.casting_time LIKE ? AND ((SUBSTRING_INDEX(basic_spells.`range`, ' ', -1) = 'feet' AND SUBSTRING_INDEX(basic_spells.`range`, ' ', 1) BETWEEN ? AND ?) OR (SUBSTRING_INDEX(basic_spells.`range`, ' ', -1) = 'mile' AND SUBSTRING_INDEX(basic_spells.`range`, ' ', 1) BETWEEN ? AND ?) OR basic_spells.`range` IN (?,?,?,?,?)) AND basic_spells.components LIKE ? AND basic_spells.duration LIKE ? AND basic_spells.description LIKE ? AND COALESCE(basic_spells.upcast, '') LIKE ?")
	// CHECK FOR MILE AND MILES !!!
	// KNOWN ISSUES:
	// id 34 not working (susbtring left side for specials)
	// mile and miles are different

	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(queryDataString.Id, queryDataString.Source, queryDataString.Level, queryDataString.School, queryDataString.IsRitual, queryDataString.CastingTime, queryDataString.RVSF, queryDataString.RVEF, queryDataString.RVSM, queryDataString.RVEM, queryDataString.Specials[0], queryDataString.Specials[1], queryDataString.Specials[2], queryDataString.Specials[3], queryDataString.Specials[4], "%", "%", "%", "%")
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

// Data integrity checks

func integrityInt(x int) string {
	if x != 0 {
		return strconv.Itoa(x)
	}
	return "%"
}

func integrityLevel(x int) string {
	if x != 0 {
		return strconv.Itoa(x - 1) // I hate that shit
	}
	return "%"
}

func integrityRangeStart(val int) int {
	if val < 0 || val > 5000 {
		return 0
	}
	return val
}

func integrityRangeEnd(val int) int {
	if val < 0 || val > 5000 {
		return 5000
	}
	return val
}

func specials(val string) [5]string {
	var (
		special = [5]string{"Self", "Sight", "Special", "Touch", "Unlimited"}
		valArr  []string
	)

	if len(val) != 5 {
		val = "11111"
	}
	valArr = strings.Split(val, "")

	for x := range valArr {
		if valArr[x] == "0" {
			special[x] = "."
		}
	}
	return special
}

func rangeValue(rangeStartFeet int, rangeStopFeet int, rangeStartMile int, rangeStopMile int) [4]int {
	//(SUBSTRING_INDEX(basic_spells.`range`, " ", -1) = "feet" AND SUBSTRING_INDEX(basic_spells.`range`, " ", 1) BETWEEN 0 AND 5000) OR (SUBSTRING_INDEX(basic_spells.`range`, " ", -1) = "mile" AND SUBSTRING_INDEX(basic_spells.`range`, " ", 1) BETWEEN 0 AND 5000) OR basic_spells.`range` IN ("","","","","") OR basic_spells.`range` LIKE ""
	// Fuck my life
	if (rangeStartFeet < 0) || (rangeStartFeet > 5000) {
		rangeStartFeet = 0
	}
	if (rangeStopFeet > 5000) || (rangeStopFeet < 0) {
		rangeStopFeet = 5000
	}
	if (rangeStartMile < 0) || (rangeStartMile > 5000) {
		rangeStartMile = 0
	}
	if (rangeStopMile > 5000) || (rangeStopMile < 0) {
		rangeStopMile = 5000
	}
	arr := [4]int{rangeStartFeet, rangeStopFeet, rangeStartMile, rangeStopMile}

	return arr
}
