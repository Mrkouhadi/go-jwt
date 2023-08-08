package main

import (
	"fmt"
	"net/http"
)

// Getjwt handler, accessing it on  "/jwt"
func GetJWT(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Api_key") == app.APIKEY {
		t, err := GetjwtToken(r)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Access Denied"))
			return
		}
		fmt.Fprint(w, t) // t is: header.payload.signature
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Access Denied"))
	}
}

// home handler, accessing it on "/"
func Home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello world"))
}
