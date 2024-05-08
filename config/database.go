package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func DatabaseConnection() (*sql.DB, error) {

	dbUsername := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	dbParams := os.Getenv("DB_PARAMS")

	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?%s",
		dbUsername,
		dbPassword,
		dbHost,
		dbPort,
		dbName,
		dbParams)

	var err error
	db, err := sql.Open("postgres", connString)
	if err != nil {
		log.Fatal("Error connecting to the database: ", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Error pinging the database: ", err)
	}

	return db, nil
}
