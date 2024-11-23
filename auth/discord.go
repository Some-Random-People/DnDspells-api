package auth

import (
	"bytes"
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
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

func createJWT(identifier string) (string, error) {
	secret := os.Getenv("SECRET")
	jwtToken := jwt.New(jwt.SigningMethodHS256)
	claims := jwtToken.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(168 * time.Hour)
	claims["identifier"] = identifier
	claims["method"] = "discord"
	stringToken, err := jwtToken.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return stringToken, nil
}

func DiscordConfig(router *mux.Router, db *sql.DB) {
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

		discordToken, err := authConf.Exchange(context.Background(), r.FormValue("code"))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			return
		}
		res, err := authConf.Client(context.Background(), discordToken).Get("https://discord.com/api/users/@me")

		if err != nil || res.StatusCode != 200 {
			w.WriteHeader(http.StatusInternalServerError)
			if err != nil {
				log.Println(err)
			} else {
				fmt.Fprint(w, res.Status)
			}
			return
		}

		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Print(err)
			return
		}

		identifier := ""
		var (
			count int
		)
		rows, err := db.Query("SELECT COUNT(*) FROM external_user_id WHERE external_method = 'discord' AND token = ?", identifier)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()
		for rows.Next() {
			if err := rows.Scan(&count); err != nil {
				log.Fatal(err)
			}
		}

		/*if count == 0 {
		    usersInsert, err := db.Prepare("INSERT INTO users(nickname)")
		    if err != nil {
		        log.Fatal(err)
		    }
		    res, err := usersInsert.Exec("")
		}*/

		jwtToken, err := createJWT(identifier)
		if err != nil {
			fmt.Println(err.Error())
		}
		cookie := http.Cookie{
			Name:     "accessToken",
			Value:    jwtToken,
			Path:     "/",
			MaxAge:   604800,
			HttpOnly: true,
			Secure:   true,
		}
		http.SetCookie(w, &cookie)

		//fmt.Fprint(w, body)
		w.Write(body)
	})
}
