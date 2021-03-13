module server/src/getty

go 1.16

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/gorilla/websocket v1.4.2
	godori.com/db v0.0.0
	godori.com/util/filter v0.0.0
	golang.org/x/oauth2 v0.0.0-20210311163135-5366d9dc1934
)

replace (
	godori.com/db v0.0.0 => ./../db
	godori.com/util/filter v0.0.0 => ./../../util/filter
)
