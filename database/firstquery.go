package database

import (
	"database/sql"
	"log"
)

func BasicSelect(db *sql.DB) {
	var (
		id   int
		name string
	)

	stmt, err := db.Prepare("SELECT id, name FROM books")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&id, &name)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(id, name)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
}
