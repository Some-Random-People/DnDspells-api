package database

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func ConnectToDatabase() *sql.DB { // This function only establishes sql.DB object, verifies it and returns it
	err := godotenv.Load() // Check for .env errors
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// Awful declaration but works
	db, err := sql.Open("mysql", os.Getenv("DB_USER")+":"+os.Getenv("DB_PASS")+"@("+os.Getenv("DB_IP")+":"+os.Getenv("DB_PORT")+")/"+os.Getenv("DB_NAME"))
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping() // Qick connection check
	if err != nil {
		log.Fatal("No connection to database")
	}

	return db // Return pointer to databse for further connections
}
