module server

go 1.16

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/gorilla/mux v1.8.0
	godori.com/db v0.0.0
	godori.com/game v0.0.0
	godori.com/getty v0.0.0
	godori.com/packet/toClient v0.0.0
	godori.com/packet/toServer v0.0.0
	godori.com/util/constant/loginType v0.0.0
	godori.com/util/constant/mapType v0.0.0
	godori.com/util/constant/modeType v0.0.0
	godori.com/util/constant/modelType v0.0.0
	godori.com/util/constant/roomType v0.0.0
	godori.com/util/constant/teamType v0.0.0
	godori.com/util/filter v0.0.0
	godori.com/util/math v0.0.0
	godori.com/util/pix v0.0.0
)

replace (
	godori.com/db v0.0.0 => ./src/db
	godori.com/game v0.0.0 => ./src/game
	godori.com/getty v0.0.0 => ./src/getty
	godori.com/packet/toClient v0.0.0 => ./src/packet/toClient
	godori.com/packet/toServer v0.0.0 => ./src/packet/toServer
	godori.com/util/constant/loginType v0.0.0 => ./util/constant/loginType
	godori.com/util/constant/mapType v0.0.0 => ./util/constant/mapType
	godori.com/util/constant/modeType v0.0.0 => ./util/constant/modeType
	godori.com/util/constant/modelType v0.0.0 => ./util/constant/modelType
	godori.com/util/constant/roomType v0.0.0 => ./util/constant/roomType
	godori.com/util/constant/teamType v0.0.0 => ./util/constant/teamType
	godori.com/util/filter v0.0.0 => ./util/filter
	godori.com/util/math v0.0.0 => ./util/math
	godori.com/util/pix v0.0.0 => ./util/pix
)
