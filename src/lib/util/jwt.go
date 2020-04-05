package util

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

func MakeSessionToken(username string, signKey []byte) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  username,
		"iat": time.Now(),
		"exp": time.Now().Add(time.Hour * 24).Unix(), // Expire in 24h
	})
	ss, err := token.SignedString(signKey)
	return ss, err
}

func ParseToken(signed string, signKey []byte) (*jwt.Token, error) {
	token, err := jwt.Parse(signed, func(token *jwt.Token) (interface{}, error) {
		return signKey, nil
	})

	return token, err
}
