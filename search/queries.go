package search

import (
	"database/sql"
	"encoding/json"
	"log"
)

func SpellList(db *sql.DB) string {
	var (
		id        int
		name      string
		spellList = make(map[string]int)
	)

	rows, err := db.Query("SELECT id, name FROM table_view ORDER BY name ASC")
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

func SchoolList(db *sql.DB) string {
	var (
		id         int
		name       string
		schoolList = make(map[string]int)
	)

	rows, err := db.Query("SELECT id, name FROM schools ORDER BY name ASC")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&id, &name)
		if err != nil {
			log.Fatal(err)
		}
		schoolList[name] = id
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	jsonBytes, err := json.Marshal(schoolList)
	if err != nil {
		log.Fatal(err)
	}

	return string(jsonBytes)
}

func BookList(db *sql.DB) string {
	var (
		id       int
		name     string
		bookList = make(map[string]int)
	)

	rows, err := db.Query("SELECT id, name FROM books ORDER BY name ASC")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&id, &name)
		if err != nil {
			log.Fatal(err)
		}
		bookList[name] = id
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	jsonBytes, err := json.Marshal(bookList)
	if err != nil {
		log.Fatal(err)
	}

	return string(jsonBytes)
}

func Misc(db *sql.DB) string {
	var (
		ct  []string
		rn  []string
		cmp []string
		dur []string
	)

	castingTime(db, &ct) // ct - Casting Time
	spellRange(db, &rn)  // rn - Range
	components(db, &cmp) // cmp - Components
	duration(db, &dur)   // dur - Duration

	type jason struct {
		CastingTime []string
		SpellRange  []string
		Components  []string
		Duration    []string
	}

	jas := &jason{
		CastingTime: ct,
		SpellRange:  rn,
		Components:  cmp,
		Duration:    dur,
	}

	jsonBytes, err := json.Marshal(jas)
	if err != nil {
		log.Fatal(err)
	}

	return string(jsonBytes)
}

func DataQuery(db *sql.DB) string {
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
		Upcast      string
	}
	var queryResult qr

	stmt, err := db.Prepare("SELECT basic_spells.id, basic_spells.`name`, books.`name`, basic_spells.`level`, schools.`name`, basic_spells.is_ritual, basic_spells.casting_time, basic_spells.`range`, basic_spells.components, basic_spells.duration, basic_spells.`description`, basic_spells.upcast FROM basic_spells JOIN books ON basic_spells.`source` LIKE books.id JOIN schools ON basic_spells.school LIKE schools.id WHERE basic_spells.id LIKE ? AND basic_spells.`source` LIKE ? AND basic_spells.level LIKE ? AND basic_spells.school LIKE ? AND basic_spells.is_ritual LIKE ? AND basic_spells.casting_time LIKE ? AND basic_spells.range LIKE ? AND basic_spells.components LIKE ? AND basic_spells.duration LIKE ? AND basic_spells.description LIKE ? AND basic_spells.upcast LIKE ?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

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
		log.Fatal(err)
	}

	return string(resultsJSON)
}

// Everything below is for "Misc" function

func spellRange(db *sql.DB, rn *[]string) {

	// Spell range query

	var temp string

	rows, err := db.Query("SELECT DISTINCT `range` FROM table_view ORDER BY `range` ASC")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&temp)
		if err != nil {
			log.Fatal(err)
		}
		*rn = append(*rn, temp)
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
}

func castingTime(db *sql.DB, ct *[]string) {

	// Casting time query

	var temp string

	rows, err := db.Query("SELECT DISTINCT casting_time FROM table_view ORDER BY casting_time ASC")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&temp)
		if err != nil {
			log.Fatal(err)
		}
		*ct = append(*ct, temp)
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
}

func components(db *sql.DB, cmp *[]string) {

	// Spell components query

	var temp string

	rows, err := db.Query("SELECT DISTINCT components FROM table_view ORDER BY components ASC")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&temp)
		if err != nil {
			log.Fatal(err)
		}
		*cmp = append(*cmp, temp)
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
}

func duration(db *sql.DB, dur *[]string) {

	// Spell duration query

	var temp string

	rows, err := db.Query("SELECT DISTINCT duration FROM table_view ORDER BY duration ASC")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&temp)
		if err != nil {
			log.Fatal(err)
		}
		*dur = append(*dur, temp)
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
}
