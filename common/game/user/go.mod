module server/common/game/user

go 1.16

require (
	godori.com/game/character v0.0.0
	godori.com/getty v0.0.0
)

replace (
	godori.com/game/character v0.0.0 => ../../game/character
	godori.com/getty v0.0.0 => ../../getty
)
