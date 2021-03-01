module server

go 1.16

require (
	github.com/gorilla/websocket v1.4.2
	godori.com/game v0.0.0
	godori.com/getty v0.0.0
	godori.com/packet/toserver v0.0.0
)

replace (
	godori.com/game v0.0.0 => ./common/game
	godori.com/getty v0.0.0 => ./common/getty
	godori.com/packet/toserver v0.0.0 => ./common/packet/toserver
)
