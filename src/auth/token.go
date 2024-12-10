package auth

import (
	"api-bk/src/config"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
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

func ValidarToken(r *http.Request) error {
	tokenString := extrairToken(r)
	token, err := jwt.Parse(tokenString, returnVerificationKey)

	if err != nil || !token.Valid {
		return err
	}

	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return nil
	}

	return errors.New("invalid token")
}

func extrairToken(r *http.Request) string {
	token := r.Header.Get("Authorization")

	if len(strings.Split(token, " ")) == 2 {
		return strings.Split(token, " ")[1]
	}

	return ""
}

func ExtrairUsuarioID(r *http.Request) (uint64, error) {
	tokenString := extrairToken(r)
	token, err := jwt.Parse(tokenString, returnVerificationKey)

	if err != nil || !token.Valid {
		return 0, err
	}

	if roles, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		usuarioID, err := strconv.ParseUint(fmt.Sprintf("%.0f", roles["usuarioID"]), 10, 64)
		if err != nil {
			return 0, err
		}

		return usuarioID, nil
	}

	return 0, errors.New("invalid token")
}

func returnVerificationKey(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("sign method unxpected! %v", token.Header["alg"])
	}

	return config.SecretKey, nil
}
