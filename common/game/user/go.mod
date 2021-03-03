module server/common/game/user

go 1.16

require (
	godori.com/database v0.0.0
	godori.com/game/character v0.0.0
	godori.com/getty v0.0.0
)

replace (
	godori.com/database v0.0.0 => ../../database
	godori.com/game/character v0.0.0 => ../../game/character
	godori.com/getty v0.0.0 => ../../getty
)
