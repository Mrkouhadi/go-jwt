package main

import (
	"net/http"
	"os"
)

/*
* In order to get authorized: you need to put the api_key in the headers on postman
 */

var app Config

func main() {
	SetUpEnvVars()                               //set up env vars
	app.APIKEY = os.Getenv("API_KEY")            // get jwt secret signature
	app.Secret = []byte(os.Getenv("JWT_SECRET")) // get jwt secret signature

	http.Handle("/", ValidateJWT(Home))
	http.HandleFunc("/jwt", GetJWT)

	http.ListenAndServe(":8080", nil)
}
