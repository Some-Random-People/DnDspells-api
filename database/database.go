package database

import (
	"database/sql"
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func ConnectToDatabase() *sql.DB { // This function only establishes sql.DB object, verifies it and returns it
	err := godotenv.Load() // Check for .env errors
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// Awful declaration but works
	//db, err := sql.Open("mysql", os.Getenv("DB_USER")+":"+os.Getenv("DB_PASS")+"@("+os.Getenv("DB_IP")+":"+os.Getenv("DB_PORT")+")/"+os.Getenv("DB_NAME"))
	cfg := mysql.Config{
		User:                 os.Getenv("DB_USER"),
		Passwd:               os.Getenv("DB_PASS"),
		Net:                  "tcp",
		Addr:                 os.Getenv("DB_IP") + ":" + os.Getenv("DB_PORT"),
		DBName:               os.Getenv("DB_NAME"),
		AllowNativePasswords: true,
	}
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping() // Qick connection check
	if err != nil {
		log.Fatal("No connection to database")
	}
	log.Println("Connected to database")

	return db // Return pointer to databse for further connections
}
