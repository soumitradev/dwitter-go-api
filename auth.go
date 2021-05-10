package main

import (
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func CreateToken(userid string) (string, error) {

	os.Setenv("ACCESS_SECRET", "MYVERYSECRETKEY")

	token_claims := jwt.MapClaims{}
	token_claims["authorized"] = true
	token_claims["username"] = userid
	token_claims["exp"] = time.Now().Add(time.Minute * 15).Unix()

	acces_token := jwt.NewWithClaims(jwt.SigningMethodHS256, token_claims)

	token, err := acces_token.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return "", err
	}

	return token, nil
}

func RefreshToken(userid string) (string, error) {

	os.Setenv("REFRESH_SECRET", "MYOTHERVERYSECRETKEY")

	token_claims := jwt.MapClaims{}
	token_claims["authorized"] = true
	token_claims["username"] = userid
	token_claims["exp"] = time.Now().Add(time.Hour * 24 * 7).Unix()

	acces_token := jwt.NewWithClaims(jwt.SigningMethodHS256, token_claims)

	token, err := acces_token.SignedString([]byte(os.Getenv("REFRESH_SECRET")))
	if err != nil {
		return "", err
	}

	return token, nil
}
