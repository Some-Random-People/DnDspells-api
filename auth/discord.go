package auth

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
)

var discordEndpoints = oauth2.Endpoint{
	AuthURL:   "https://discord.com/api/oauth2/authorize",
	TokenURL:  "https://discord.com/api/oauth2/token",
	AuthStyle: oauth2.AuthStyleInParams,
}

func pad(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}
func generateStateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func encryptToken(token []byte) (string, error) {
	secret := os.Getenv("SECRET")
	block, err := aes.NewCipher([]byte(secret))
	if err != nil {
		return "", err
	}
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}
	cbc := cipher.NewCBCEncrypter(block, iv)

	paddedToken := pad(token, aes.BlockSize)

	encrypted := make([]byte, len(paddedToken))

	cbc.CryptBlocks(encrypted, paddedToken)
	return hex.EncodeToString(append(iv, encrypted...)), nil
}

func DiscordConfig(router *mux.Router) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error while loading .env file")
	}
	var store = sessions.NewCookieStore([]byte(os.Getenv("COOKIE")))
	authConf := &oauth2.Config{
		RedirectURL:  "http://127.0.0.1:80/api/auth/redirect",
		ClientID:     os.Getenv("DC_CLIENTID"),
		ClientSecret: os.Getenv("DC_SECRET"),
		Scopes:       []string{"identify", "email"},
		Endpoint:     discordEndpoints,
	}

	router.HandleFunc("/api/auth/discord", func(w http.ResponseWriter, r *http.Request) {
		state, err := generateStateToken()
		if err != nil {
			log.Fatal("Error with status token generation!")
		}
		session, _ := store.Get(r, "auth-session")
		session.Values["state"] = state
		session.Save(r, w)
		http.Redirect(w, r, authConf.AuthCodeURL(state), http.StatusTemporaryRedirect)
	})

	router.HandleFunc("/api/auth/redirect", func(w http.ResponseWriter, r *http.Request) {
		authSession, _ := store.Get(r, "auth-session")
		if r.FormValue("state") != authSession.Values["state"] {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		token, err := authConf.Exchange(context.Background(), r.FormValue("code"))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			return
		}
		res, err := authConf.Client(context.Background(), token).Get("https://discord.com/api/users/@me")

		if err != nil || res.StatusCode != 200 {
			w.WriteHeader(http.StatusInternalServerError)
			if err != nil {
				log.Println(err)
			} else {
				fmt.Fprint(w, res.Status)
			}
			return
		}
		tokenJSON, err := json.Marshal(token)
		if err != nil {
			log.Println(err)
			return
		}
		accessToken, err := encryptToken(tokenJSON)
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Println(accessToken)
		cookie := http.Cookie{
			Name:     "accessToken",
			Value:    accessToken,
			Path:     "/",
			MaxAge:   604800,
			HttpOnly: true,
			Secure:   true,
		}
		http.SetCookie(w, &cookie)

		fmt.Fprint(w, "Registered")
	})
}
