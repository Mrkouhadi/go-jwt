package main

import (
	"fmt"
	"net/http"
)

// Getjwt handler, accessing it on  "/jwt": will outpu a token ==> header.payload.signature
func GetJWT(w http.ResponseWriter, r *http.Request) {

	/**** this supposed to be done by the user client
	r.Header.Set("Api_key", "123456789")
	****/

	if r.Header.Get("Api_key") == app.APIKEY {
		t, err := GetjwtToken(r)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Access Denied"))
			return
		}
		fmt.Fprint(w, t) // ==> t is: header.payload.signature
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Access Denied"))
	}
}

// home handler, accessing it on "/", will result in access denied
func Home(w http.ResponseWriter, r *http.Request) {
	// GetJWT gives us a token and when the client takes that and put it in the header as  "Token" which the key...
	// then the client will see the message below
	w.Write([]byte("Hello world"))
}
