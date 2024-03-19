package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
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
			redirectURL := fmt.Sprintf("/refresh-token?redirect=%s", url.QueryEscape(r.URL.Path))
			http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
			return
		}
		//// Access claims from the token
		claims, ok := validToken.Claims.(*Claims)
		if !ok {
			ErrorJSON(w, err, http.StatusUnauthorized)
			return
		}
		//// Check if the token is expired
		expiresAtUnix := claims.ExpiresAt.Time.Unix()
		if time.Now().Unix() >= expiresAtUnix {
			redirectURL := fmt.Sprintf("/refresh-token?redirect=%s", url.QueryEscape(r.URL.Path))
			http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
			return
		}
		next.ServeHTTP(w, r)
	})
}
