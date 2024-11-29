package auth

import (
	"api-bk/src/config"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

func CriarToken(usuarioID uint64) (string, error) {
	roles := jwt.MapClaims{}

	roles["authorized"] = true
	roles["exp"] = time.Now().Add(time.Hour * 6).Unix()
	roles["usuarioID"] = usuarioID

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, roles)

	return token.SignedString([]byte(config.SecretKey))
}
