package auth

import (
	"crypto/rand"
	"log"
	"math/big"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
)

var discordEndpoints = oauth2.Endpoint{
	AuthURL:   "https://discord.com/api/oauth2/authorize",
	TokenURL:  "https://discord.com/api/oauth2/token",
	AuthStyle: oauth2.AuthStyleInParams,
}

func generateStateToken() string {
	token, err := rand.Int(rand.Reader, big.NewInt(100000))
	if err != nil {
		log.Fatal("Failed to create state token!")
	}
	return token.String()
}
func DiscordConfig(router *mux.Router) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error while loading .env file")
	}
	authConf := &oauth2.Config{
		RedirectURL:  "http://127.0.0.1:80/auth/redirect",
		ClientID:     os.Getenv("DC_CLIENTID"),
		ClientSecret: os.Getenv("DC_SECRET"),
		Scopes:       []string{"identify", "email"},
		Endpoint:     discordEndpoints,
	}

	router.HandleFunc("/auth/discord", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, authConf.AuthCodeURL(generateStateToken()), http.StatusTemporaryRedirect)
	})

	router.HandleFunc("/auth/redirect", func(w http.ResponseWriter, r *http.Request) {
	})
}
