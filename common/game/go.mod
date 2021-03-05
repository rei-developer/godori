module server/common/game/user

go 1.16

require (
	godori.com/db v0.0.0
	godori.com/getty v0.0.0
	godori.com/util/math v0.0.0
	godori.com/util/pix v0.0.0
)

replace (
	godori.com/db v0.0.0 => ./../db
	godori.com/getty v0.0.0 => ./../getty
	godori.com/util/math v0.0.0 => ./../../util/math
	godori.com/util/pix v0.0.0 => ./../../util/pix
)
