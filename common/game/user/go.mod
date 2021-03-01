module server/common/game/user

go 1.16

require (
	godori.com/game v0.0.0
	godori.com/getty v0.0.0
)

replace (
	godori.com/game v0.0.0 => ../../../common/game
	godori.com/getty v0.0.0 => ../../../common/getty
)
