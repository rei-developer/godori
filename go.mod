module server

go 1.16

require (
	godori.com/database v0.0.0
	godori.com/game v0.0.0
	godori.com/getty v0.0.0
	godori.com/packet/toServer v0.0.0
)

replace (
	godori.com/database v0.0.0 => ./common/database
	godori.com/game v0.0.0 => ./common/game
	godori.com/getty v0.0.0 => ./common/getty
	godori.com/packet/toServer v0.0.0 => ./common/packet/toServer
	godori.com/util/math v0.0.0 => ./util/math
	godori.com/util/pix v0.0.0 => ./util/pix
)
