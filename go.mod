module server

go 1.16

require (
	github.com/go-sql-driver/mysql v1.5.0 // indirect
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/joho/godotenv v1.3.0 // indirect
	godori.com/database v0.0.0 // indirect
	godori.com/game/character v0.0.0 // indirect
	godori.com/game/shop v0.0.0
	godori.com/game/user v0.0.0
	godori.com/getty v0.0.0
	godori.com/packet/toserver v0.0.0
)

replace (
	godori.com/database v0.0.0 => ./common/database
	godori.com/game/character v0.0.0 => ./common/game/character
	godori.com/game/shop v0.0.0 => ./common/game/shop
	godori.com/game/user v0.0.0 => ./common/game/user
	godori.com/getty v0.0.0 => ./common/getty
	godori.com/packet/toserver v0.0.0 => ./common/packet/toserver
)
