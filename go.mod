module server

go 1.16

require (
	github.com/gorilla/websocket v1.4.2
	sandspoon.com/game v0.0.0
	sandspoon.com/getty v0.0.0
	sandspoon.com/packet v0.0.0
)

replace (
	sandspoon.com/game v0.0.0 => ./common/game
	sandspoon.com/getty v0.0.0 => ./common/getty
	sandspoon.com/packet v0.0.0 => ./common/packet
)
