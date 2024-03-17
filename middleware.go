package main

import (
	"errors"
	"log"
	"net/http"
	"time"
)

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		secKey, err := getSecKey()
		if err != nil {
			log.Fatal(err)
		}
		value, err := ReadEncrypted(r, "access_token", secKey)
		if err != nil {
			ErrorJSON(w, err, http.StatusUnauthorized)
			return
		}
		// Validate the access token
		validToken, err := validateToken(value)
		if err != nil || !validToken.Valid {
			ErrorJSON(w, err, http.StatusUnauthorized)
			return
		}
		// Access claims from the token
		claims, ok := validToken.Claims.(*Claims)
		if !ok {
			ErrorJSON(w, err, http.StatusUnauthorized)
			return
		}
		// Check if the access token is expired
		expiresAtUnix := claims.ExpiresAt.Time.Unix()
		if time.Now().Unix() >= expiresAtUnix {
			ErrorJSON(w, err, http.StatusUnauthorized)
			return
		}
		// we check first if we have them in Redis
		creds, err := ReadCredentials(client, claims.Email)
		if err != nil {
			log.Println("Error reading credentials from Redis:", err)
			ErrorJSON(w, err, http.StatusUnauthorized)
			return
		}

		// Compare retrieved credentials with claims
		if creds != claims.Password {
			log.Println("Invalid credentials")
			ErrorJSON(w, errors.New("unauthorized"), http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
