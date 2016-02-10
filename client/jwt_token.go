package main

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

func generateToken() string {
	token := jwt.New(jwt.SigningMethodHS256)

	// Set a header and a claim
	token.Header["typ"] = "JWT"
	token.Claims["exp"] = time.Now().Add(time.Hour * 96).Unix()

	// Generate encoded token
	t, err := token.SignedString([]byte(SigningKey))
	checkError(err)
	return t
}
