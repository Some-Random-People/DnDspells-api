module github.com/some-random-people/dndspells-api

go 1.23.3

require github.com/gorilla/mux v1.8.1 // direct

require (
	github.com/go-sql-driver/mysql v1.8.1
	github.com/golang-jwt/jwt/v5 v5.2.2
	github.com/gorilla/sessions v1.4.0
	github.com/joho/godotenv v1.5.1
	golang.org/x/oauth2 v0.24.0
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/gorilla/securecookie v1.1.2 // indirect
)
