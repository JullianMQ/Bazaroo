package db

import (
	"database/sql"
	"log"
	"net/url"
	"os"

	"github.com/joho/godotenv"
)

var DB *sql.DB

func ConnDB() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	serviceURI := os.Getenv("AIVEN_DB_URI")
	conn, _ := url.Parse(serviceURI)
	conn.RawQuery = "sslmode=verify-ca;sslrootcert=ca.pem"

	DB, err = sql.Open("postgres", conn.String())

	if err != nil {
		log.Fatal(err)
	}
}

