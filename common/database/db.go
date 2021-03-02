package database

import (
	"database/sql"
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
	checkError(err)
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
	checkError(err)
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	err = db.Ping() // This DOES open a connection if necessary. This makes sure the database is accessible
	checkError(err)
}

func GetUserCount() (count int) {
	err := db.QueryRow("SELECT COUNT(*) count FROM users").Scan(&count)
	checkError(err)
	return count
}

func GetUser(findIndex int) (id string, uuid string) {
	err := db.QueryRow("SELECT id, uuid FROM users WHERE `index` = ?", findIndex).Scan(&id, &uuid)
	checkError(err)
	return id, uuid
}

type UserData struct {
	Id   string
	Uuid string
	Name string
}

func GetUsers() []UserData {
	rows, err := db.Query("SELECT id, uuid, name FROM users")
	checkError(err)
	defer rows.Close()
	var items []UserData
	for rows.Next() {
		item := UserData{}
		err = rows.Scan(
			&item.Id,
			&item.Uuid,
			&item.Name,
		)
		checkError(err)
		items = append(items, item)
	}
	return items
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
