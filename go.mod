module server

go 1.16

require (
	godori.com/db v0.0.0
	godori.com/game v0.0.0
	godori.com/getty v0.0.0
	godori.com/packet/toClient v0.0.0
	godori.com/packet/toServer v0.0.0
	godori.com/util/constant/modeType v0.0.0
	godori.com/util/constant/roomType v0.0.0
	godori.com/util/math v0.0.0
	godori.com/util/pix v0.0.0
)

replace (
	godori.com/db v0.0.0 => ./src/db
	godori.com/game v0.0.0 => ./src/game
	godori.com/getty v0.0.0 => ./src/getty
	godori.com/packet/toClient v0.0.0 => ./src/packet/toClient
	godori.com/packet/toServer v0.0.0 => ./src/packet/toServer
	godori.com/util/constant/modeType v0.0.0 => ./util/constant/modeType
	godori.com/util/constant/roomType v0.0.0 => ./util/constant/roomType
	godori.com/util/math v0.0.0 => ./util/math
	godori.com/util/pix v0.0.0 => ./util/pix
)