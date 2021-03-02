module server/common/getty

go 1.16

require (
	github.com/gorilla/websocket v1.4.2
	godori.com/database v0.0.0
)

replace (
	godori.com/database v0.0.0 => ../../common/database
)