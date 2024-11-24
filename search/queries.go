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

	castingTime(db, &ct)
	spellRange(db, &rn)
	components(db, &cmp)
	duration(db, &dur)

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
