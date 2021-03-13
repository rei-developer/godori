package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var database *sql.DB

type Portal struct {
	Place     sql.NullInt32
	X         sql.NullInt32
	Y         sql.NullInt32
	NextPlace sql.NullInt32
	NextX     sql.NullInt32
	NextY     sql.NullInt32
	NextDirX  sql.NullInt32
	NextDirY  sql.NullInt32
	Sound     sql.NullString
}

type Item struct {
	Id          sql.NullInt32
	Num         sql.NullInt32
	Icon        sql.NullString
	Name        sql.NullString
	Description sql.NullString
	Cost        sql.NullInt32
	Method      sql.NullString
}

type Inventory struct {
	ItemId sql.NullInt32
	Num    sql.NullInt32
	Expiry sql.NullInt32
}

type Billing struct {
	Id                 sql.NullInt32
	ProductId          sql.NullInt32
	PurchaseDate       sql.NullTime
	UseState           sql.NullInt32
	RefundRequestState sql.NullInt32
}

type NoticeMessage struct {
	Id      sql.NullInt32
	Title   sql.NullString
	Created sql.NullTime
	Deleted sql.NullTime
	Avatar  sql.NullString
	Author  sql.NullString
}

type Rank struct {
	Id       sql.NullInt32
	Name     sql.NullString
	Level    sql.NullInt32
	Exp      sql.NullInt32
	Point    sql.NullInt32
	Kill     sql.NullInt32
	Death    sql.NullInt32
	Assist   sql.NullInt32
	Avatar   sql.NullString
	Memo     sql.NullString
	Admin    sql.NullInt32
	Clanname sql.NullString
}

type Clan struct {
	Id         sql.NullInt32
	MasterId   sql.NullInt32
	Name       sql.NullString
	Level1Name sql.NullString
	Level2Name sql.NullString
	Level3Name sql.NullString
	Level4Name sql.NullString
	Level5Name sql.NullString
	Notice     sql.NullString
	Level      sql.NullInt32
	Exp        sql.NullInt32
	Cash       sql.NullInt32
	Coin       sql.NullInt32
	Regdate    sql.NullTime
	Condition  sql.NullInt32
}

type User struct {
	Id           sql.NullInt32
	Uid          sql.NullString
	Uuid         sql.NullString
	Name         sql.NullString
	Sex          sql.NullInt32
	Level        sql.NullInt32
	Exp          sql.NullInt32
	Coin         sql.NullInt32
	Cash         sql.NullInt32
	Point        sql.NullInt32
	Win          sql.NullInt32
	Lose         sql.NullInt32
	Kill         sql.NullInt32
	Death        sql.NullInt32
	Assist       sql.NullInt32
	Blast        sql.NullInt32
	Rescue       sql.NullInt32
	Survive      sql.NullInt32
	Escape       sql.NullInt32
	RedGraphics  sql.NullString
	BlueGraphics sql.NullString
	Memo         sql.NullString
	LastChat     sql.NullInt32
	Admin        sql.NullInt32
	Verify       sql.NullInt32
}

func init() {
	var err error
	driverName, dataSourceName, connectionLimit := LoadConfig()
	database, err = sql.Open(driverName, dataSourceName)
	CheckError(err)
	database.SetConnMaxLifetime(time.Minute)
	database.SetMaxOpenConns(connectionLimit)
	database.SetMaxIdleConns(connectionLimit)
	err = database.Ping()
	CheckError(err)
}

func LoadConfig() (driverName string, dataSourceName string, connectionLimit int) {
	err := godotenv.Load()
	CheckError(err)
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

func GetPortals(place int) []Portal {
	rows, err := database.Query("SELECT place, x, y, next_place, next_x, next_y, next_dir_x, next_dir_y, sound FROM portals WHERE place = ?", place)
	CheckError(err)
	defer rows.Close()
	items := []Portal{}
	for rows.Next() {
		item := Portal{}
		err = rows.Scan(
			&item.Place,
			&item.X,
			&item.Y,
			&item.NextPlace,
			&item.NextX,
			&item.NextY,
			&item.NextDirX,
			&item.NextDirY,
			&item.Sound,
		)
		CheckError(err)
		items = append(items, item)
	}
	return items
}

func GetItems() []Item {
	rows, err := database.Query("SELECT id, num, icon, name, description, cost, method FROM items")
	CheckError(err)
	defer rows.Close()
	items := []Item{}
	for rows.Next() {
		item := Item{}
		err = rows.Scan(
			&item.Id,
			&item.Num,
			&item.Icon,
			&item.Name,
			&item.Description,
			&item.Cost,
			&item.Method,
		)
		CheckError(err)
		items = append(items, item)
	}
	return items
}

func GetInventorys() []Inventory {
	rows, err := database.Query("SELECT item_id, num, expiry FROM inventorys")
	CheckError(err)
	defer rows.Close()
	items := []Inventory{}
	for rows.Next() {
		item := Inventory{}
		err = rows.Scan(
			&item.ItemId,
			&item.Num,
			&item.Expiry,
		)
		CheckError(err)
		items = append(items, item)
	}
	return items
}

func GetBillings() []Billing {
	rows, err := database.Query("SELECT id, productId, purchaseDate, useState, refundRequestState FROM billings")
	CheckError(err)
	defer rows.Close()
	items := []Billing{}
	for rows.Next() {
		item := Billing{}
		err = rows.Scan(
			&item.Id,
			&item.ProductId,
			&item.PurchaseDate,
			&item.UseState,
			&item.RefundRequestState,
		)
		CheckError(err)
		items = append(items, item)
	}
	return items
}

func GetNoticeMessages(id int, deleted bool) []NoticeMessage {
	isDeleted := ""
	if deleted {
		isDeleted = "!"
	}
	rows, err := database.Query(`
		SELECT
			nm.id,
			nm.title,
			nm.created,
			nm.deleted,
			u.name author,
			u.blue_graphics avatar
		FROM notice_messages nm
			LEFT JOIN users u ON u.id = nm.target_id
		WHERE nm.user_id = ? AND ?ISNULL(nm.deleted)
		ORDER BY id DESC
	`, id, isDeleted)
	CheckError(err)
	defer rows.Close()
	items := []NoticeMessage{}
	for rows.Next() {
		item := NoticeMessage{}
		err = rows.Scan(
			&item.Id,
			&item.Title,
			&item.Created,
			&item.Deleted,
			&item.Author,
			&item.Avatar,
		)
		CheckError(err)
		items = append(items, item)
	}
	return items
}

func GetRanks() []Rank {
	rows, err := database.Query(`
		SELECT
			u.id,
			u.name,
			u.level,
			u.exp,
			u.point,
			u.kill,
			u.death,
			u.assist,
			u.blue_graphics avatar,
			u.memo,
			u.admin,
			c.name clanname
		FROM users u
			LEFT JOIN clan_members cm ON cm.user_id = u.id
			LEFT JOIN clans c ON c.id = cm.clan_id
		WHERE u.verify = 1 ORDER BY u.point DESC
	`)
	CheckError(err)
	defer rows.Close()
	items := []Rank{}
	for rows.Next() {
		item := Rank{}
		err = rows.Scan(
			&item.Id,
			&item.Name,
			&item.Level,
			&item.Exp,
			&item.Point,
			&item.Kill,
			&item.Death,
			&item.Assist,
			&item.Avatar,
			&item.Memo,
			&item.Admin,
			&item.Clanname,
		)
		CheckError(err)
		items = append(items, item)
	}
	return items
}

func GetClans() []Clan {
	rows, err := database.Query(`
		SELECT
			id,
		    master_id,
		    name,
		    level1_name,
		    level2_name,
		    level3_name,
		    level4_name,
		    level5_name,
		    notice,
		    level,
		    exp,
			cash,
			coin,
			regdate,
			condition
		FROM clans
	`)
	CheckError(err)
	defer rows.Close()
	items := []Clan{}
	for rows.Next() {
		item := Clan{}
		err = rows.Scan(
			&item.Id,
			&item.MasterId,
			&item.Name,
			&item.Level1Name,
			&item.Level2Name,
			&item.Level3Name,
			&item.Level4Name,
			&item.Level5Name,
			&item.Notice,
			&item.Level,
			&item.Exp,
			&item.Cash,
			&item.Coin,
			&item.Regdate,
			&item.Condition,
		)
		CheckError(err)
		items = append(items, item)
	}
	return items
}

func GetClanMembers(clanId int) (userId []int) {
	rows, err := database.Query("SELECT user_id userId FROM clan_members WHERE clan_id = ? ORDER BY level DESC", clanId)
	CheckError(err)
	defer rows.Close()
	items := []int{}
	for rows.Next() {
		userId := 0
		err = rows.Scan(&userId)
		CheckError(err)
		items = append(items, userId)
	}
	return items
}

func GetInviteClans(userId int) (clanId []int) {
	rows, err := database.Query("SELECT clan_id clanId FROM invite_clans WHERE target_id = ?", userId)
	CheckError(err)
	defer rows.Close()
	items := []int{}
	for rows.Next() {
		clanId := 0
		err = rows.Scan(&clanId)
		CheckError(err)
		items = append(items, clanId)
	}
	return items
}

func GetUser(args map[string]interface{}) (User, bool) {
	var keys []string
	var values []interface{}
	for k, v := range args {
		keys = append(keys, fmt.Sprintf("%s = ?", k))
		values = append(values, v)
	}
	cond := strings.Join(keys, " AND ")
	item := User{}
	err := database.QueryRow(`
		SELECT
			id,
			uid,
			uuid,
			name,
			sex,
			level,
			exp,
			coin,
			cash,
			point,
			win,
			lose,
			death,
			assist,
			blast,
			1,
			survive,
			escape,
			red_graphics,
			blue_graphics,
			memo,
			unix_timestamp(last_chat) last_chat,
			admin,
			verify
		FROM users
		WHERE
	`+cond, values...).Scan(
		&item.Id,
		&item.Uid,
		&item.Uuid,
		&item.Name,
		&item.Sex,
		&item.Level,
		&item.Exp,
		&item.Coin,
		&item.Cash,
		&item.Point,
		&item.Win,
		&item.Lose,
		&item.Death,
		&item.Assist,
		&item.Blast,
		&item.Rescue,
		&item.Survive,
		&item.Escape,
		&item.RedGraphics,
		&item.BlueGraphics,
		&item.Memo,
		&item.LastChat,
		&item.Admin,
		&item.Verify,
	)
	ok := CheckError(err)
	return item, ok
}

func GetUserById(id int) (User, bool) {
	var array map[string]interface{} = make(map[string]interface{})
	array["id"] = id
	return GetUser(array)
}

func GetUserByOAuth(uid string, loginType int) (User, bool) {
	var array map[string]interface{} = make(map[string]interface{})
	array["uid"] = uid
	array["login_type"] = loginType
	return GetUser(array)
}

func GetUserByName(name string) (User, bool) {
	var array map[string]interface{} = make(map[string]interface{})
	array["name"] = name
	return GetUser(array)
}

func GetUserCount() (count int) {
	err := database.QueryRow("SELECT COUNT(*) count FROM users").Scan(&count)
	CheckError(err)
	return count
}

func UpdateUserVerify(name string, uid string, lType int) {
	_, err := database.Exec("UPDATE users SET `name` = ?, verify = 1 WHERE uid = ? AND login_type = ?", name, uid, lType)
	CheckError(err)
}

func InsertUser(uid string, lType int) {
	_, err := database.Exec("INSERT INTO users (uid, login_type) VALUES (?, ?)", uid, lType)
	CheckError(err)
}

func CheckError(err error) bool {
	if err != nil {
		if err == sql.ErrNoRows {
			return false
		} else {
			log.Println(err)
		}
	}
	return true
}
