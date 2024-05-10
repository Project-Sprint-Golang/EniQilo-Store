package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() {
	godotenv.Load()
	dbUsername := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	dbParams := os.Getenv("DB_PARAMS")

	connString := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?%s",
		dbUsername,
		dbPassword,
		dbHost,
		dbPort,
		dbName,
		dbParams)

	var err error
	DB, err = sql.Open("postgres", connString)
	if err != nil {
		log.Fatal("Error connecting to the database: ", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal("Error pinging the database: ", err)
	}
	fmt.Println("Success Connect")
	// return DB, err
}
