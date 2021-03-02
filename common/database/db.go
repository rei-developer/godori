package database

import (
	"database/sql"
	"fmt"
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
	driver   string
	host     string
	port     string
	user     string
	password string
	database string
}

type UserData struct {
	Id   string
	Uuid string
	Name string
}

func init() {
	LoadConfig()
	var err error
	dataSource := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbInfo.user, dbInfo.password, dbInfo.host, dbInfo.port, dbInfo.database)
	db, err = sql.Open(dbInfo.driver, dataSource)
	checkError(err)
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	err = db.Ping()
	checkError(err)
}

func LoadConfig() {
	err := godotenv.Load()
	checkError(err)
	dbInfo = &Database{
		os.Getenv("DB_DRIVER"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_DATABASE"),
	}
}

func GetUser(findIndex int) (id string, uuid string) {
	err := db.QueryRow("SELECT id, uuid FROM users WHERE `index` = ?", findIndex).Scan(&id, &uuid)
	checkError(err)
	return id, uuid
}

func GetUsers() []UserData {
	rows, err := db.Query("SELECT id, uuid, name FROM users")
	checkError(err)
	defer rows.Close()
	items := []UserData{}
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

func GetUserCount() (count int) {
	err := db.QueryRow("SELECT COUNT(*) count FROM users").Scan(&count)
	checkError(err)
	return count
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
