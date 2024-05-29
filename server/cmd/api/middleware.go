package main

import (
	"fmt"
	"go-jwt/internal/auth"
	"go-jwt/internal/helpers"
	"net/http"
	"net/url"
	"time"
)

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//// read the cookies data
		value, err := auth.ReadEncrypted(r, "access_token", app.Auth.AesSecretKey)
		if err != nil {
			helpers.ErrorJSON(w, err, http.StatusUnauthorized)
			return
		}
		//// Validate the access token
		validToken, err := auth.ValidateToken(value)
		if err != nil || !validToken.Valid {
			redirectURL := fmt.Sprintf("/refresh-token?redirect=%s", url.QueryEscape(r.URL.Path))
			http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
			return
		}
		//// Access claims from the token
		claims, ok := validToken.Claims.(*auth.Claims)
		if !ok {
			helpers.ErrorJSON(w, err, http.StatusUnauthorized)
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
