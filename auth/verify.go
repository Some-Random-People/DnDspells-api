package auth

import (
	"fmt"
	"log"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

func VerifyToken(tokenString string) (bool, jwt.MapClaims) {
	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:]
	} else {
		return false, nil
	}
	secret := []byte(os.Getenv("SECRET"))
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("Unexpected signing method in JWT token")
		}
		return secret, nil
	})

	if err != nil {
		log.Printf("Token validation failed: %v\n", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		return true, claims
	} else {
		return false, nil
	}
}
