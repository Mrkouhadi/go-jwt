package main

import (
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

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
				return app.Secret, nil
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
