package database

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var db *sql.DB

type User struct {
	Id   sql.NullInt32
	Uid  sql.NullString
	Uuid sql.NullString
	Name sql.NullString
}

func init() {
	var err error
	driverName, dataSourceName, connectionLimit := LoadConfig()
	db, err = sql.Open(driverName, dataSourceName)
	checkError(err)
	db.SetConnMaxLifetime(time.Minute)
	db.SetMaxOpenConns(connectionLimit)
	db.SetMaxIdleConns(connectionLimit)
	err = db.Ping()
	checkError(err)
}

func LoadConfig() (driverName string, dataSourceName string, connectionLimit int) {
	err := godotenv.Load()
	checkError(err)
	var (
		driver   = os.Getenv("DB_DRIVER")
		host     = os.Getenv("DB_HOST")
		port     = os.Getenv("DB_PORT")
		user     = os.Getenv("DB_USER")
		password = os.Getenv("DB_PASSWORD")
		database = os.Getenv("DB_DATABASE")
	)
	driverName = driver
	dataSourceName = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, host, port, database)
	connectionLimit, _ = strconv.Atoi(os.Getenv("DB_CONNECTION_LIMIT"))
	return
}

func GetUser(findId int) (uid sql.NullString, uuid sql.NullString) {
	err := db.QueryRow("SELECT uid, uuid FROM users WHERE id = ?", findId).Scan(&uid, &uuid)
	checkError(err)
	return uid, uuid
}

func GetUsers() []User {
	rows, err := db.Query("SELECT id, uid, uuid, name FROM users")
	checkError(err)
	defer rows.Close()
	items := []User{}
	for rows.Next() {
		item := User{}
		err = rows.Scan(
			&item.Id,
			&item.Uid,
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
