module server/src/game/user

go 1.16

require (
	godori.com/db v0.0.0
	godori.com/getty v0.0.0
	godori.com/packet/toClient v0.0.0
	godori.com/util/constant/mapType v0.0.0
	godori.com/util/constant/modeType v0.0.0
	godori.com/util/constant/roomType v0.0.0
	godori.com/util/constant/teamType v0.0.0
	godori.com/util/math v0.0.0
	godori.com/util/pix v0.0.0
)

replace (
	godori.com/db v0.0.0 => ./../db
	godori.com/getty v0.0.0 => ./../getty
	godori.com/packet/toClient v0.0.0 => ./../packet/toClient
	godori.com/util/constant/mapType v0.0.0 => ./../../util/constant/mapType
	godori.com/util/constant/modeType v0.0.0 => ./../../util/constant/modeType
	godori.com/util/constant/roomType v0.0.0 => ./../../util/constant/roomType
	godori.com/util/constant/teamType v0.0.0 => ./../../util/constant/teamType
	godori.com/util/math v0.0.0 => ./../../util/math
	godori.com/util/pix v0.0.0 => ./../../util/pix
)
