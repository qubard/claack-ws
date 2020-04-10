package util

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

func MakeSessionToken(username string, signKey []byte) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"iat":      time.Now(),
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // Expire in 24h
	})
	// Don't think this ever errors
	// https://godoc.org/github.com/dgrijalva/jwt-go#Token.SignedString
	ss, err := token.SignedString(signKey)
	return ss, err
}

func ParseToken(signed string, signKey []byte) (*jwt.Token, error) {
	token, err := jwt.Parse(signed, func(token *jwt.Token) (interface{}, error) {
		return signKey, nil
	})

	return token, err
}

func ExtractField(token string, field string, signKey []byte) (interface{}, bool) {
	if token, err := ParseToken(token, signKey); err == nil && token.Valid {
		claims := token.Claims.(jwt.MapClaims)
		if res, ok := claims[field]; ok {
			return res, true
		}
	}
	return nil, false
}
