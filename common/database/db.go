package database

import (
	"database/sql"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var (
	db     *sql.DB
	dbInfo *Database
)

type Database struct {
	user     string
	password string
	url      string
	engine   string
	database string
}

func LoadEnvFile() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	dbInfo = &Database{
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_URL"),
		os.Getenv("DB_ENGINE"),
		os.Getenv("DB_DATABASE"),
	}
}

func init() {
	LoadEnvFile()
	var err error
	db, err = sql.Open(dbInfo.engine, dbInfo.user+":"+dbInfo.password+"@tcp("+dbInfo.url+")/"+dbInfo.database)
	if err != nil {
		log.Fatalf("Error on initializing database connection: %s", err.Error())
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	err = db.Ping() // This DOES open a connection if necessary. This makes sure the database is accessible
	if err != nil {
		log.Fatalf("Error on opening database connection: %s", err.Error())
	}
}

func GetUserCount() (count int) {
	err := db.QueryRow("SELECT COUNT(*) count FROM users").Scan(&count)
	if err != nil {
		log.Fatal(err)
	}
	return count
}
