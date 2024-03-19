package main

import (
	"errors"
	"log"
	"net/http"
	"time"
)

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//// read the cookies data
		secKey, err := getSecKey()
		if err != nil {
			log.Fatal(err)
		}
		value, err := ReadEncrypted(r, "access_token", secKey)
		if err != nil {
			ErrorJSON(w, err, http.StatusUnauthorized)
			return
		}
		//// Validate the access token
		validToken, err := validateToken(value)
		if err != nil || !validToken.Valid {
			ErrorJSON(w, err, http.StatusUnauthorized)
			return
		}
		//// Access claims from the token
		claims, ok := validToken.Claims.(*Claims)
		if !ok {
			ErrorJSON(w, err, http.StatusUnauthorized)
			return
		}
		//// Check if the access token is expired
		expiresAtUnix := claims.ExpiresAt.Time.Unix()
		if time.Now().Unix() >= expiresAtUnix {
			ErrorJSON(w, errors.New("access token is expired"), http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
