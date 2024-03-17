package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

// FIXME:validation of credentials
func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		secKey, err := getSecKey()
		if err != nil {
			log.Fatal(err)
		}
		value, err := ReadEncrypted(r, "access_token", secKey)
		if err != nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			// http.Error(w, "please log in first.", http.StatusBadRequest)
			return
		}
		// Validate the refresh token
		validToken, err := validateToken(value)
		if err != nil {
			// http.Error(w, "middleware error: Invalid refresh token", http.StatusUnauthorized)
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		// Access claims from the token
		claims, ok := validToken.Claims.(*Claims)
		if !ok {
			// http.Error(w, "Error getting claims from token", http.StatusInternalServerError)
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		// Check if the refresh token is expired
		if time.Until(time.Unix(claims.ExpiresAt, 0)) <= 0 {
			// http.Error(w, "Access token has been expired", http.StatusUnauthorized)
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		// FIXME: add more verification of the credentials; Database queries.
		email := claims.Email
		password := claims.Password
		fmt.Println(email, "::", password)
		//
		next.ServeHTTP(w, r)
	})
}

//////
/*

// To validate the jwt token
func ValidateJWT(next func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Token"] != nil {
			token, err := jwt.Parse(r.Header["Token"][0], func(t *jwt.Token) (interface{}, error) {
				_, ok := t.Method.(*jwt.SigningMethodHMAC)
				if !ok {
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte("Access Denied"))
				}
				return secretKey, nil
			})
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Access Denied: " + err.Error()))
			}
			if token.Valid {
				//this just to show that we can get the custom value added to claims in CreateJWT
				claims := token.Claims.(jwt.MapClaims) // output: accountNumber of the user
				fmt.Println("user id : ", claims["userID"])
				// end getting the custom claims value
				next(w, r)
			}
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Access Denied"))
		}
	})
}


*/
