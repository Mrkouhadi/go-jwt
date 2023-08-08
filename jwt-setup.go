package main

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// SET UP a jwt
func CreateJWT() (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims) //claims is the second part of the token
	claims["exp"] = time.Now().Add(time.Hour).Unix()
	//*** just to show u that we're able to add custom var to claims
	claims["userID"] = "hardcoded-ID"
	//*** end of adding a custom value to claims, now we can retrieve it in ValidateJWT middleware

	// GET token string
	tokenStr, err := token.SignedString(app.Secret)
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}
	return tokenStr, nil
}

func GetjwtToken(r *http.Request) (string, error) {
	if r.Header["Api_key"] != nil {
		if r.Header["Api_key"][0] == app.APIKEY {
			token, err := CreateJWT()
			if err != nil {
				return "", err
			}
			return token, nil
		}
	}
	return "", errors.New("invalid token")
}
