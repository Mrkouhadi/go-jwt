package main

import (
	"fmt"
	"net/http"
	"os"
)

/*
* In order to get authorized: you need to put the Api_key in the headers on postman/or any client...
 */

var app Config

func main() {
	SetUpEnvVars()                               //set up env vars
	app.APIKEY = os.Getenv("API_KEY")            // get jwt secret signature
	app.Secret = []byte(os.Getenv("JWT_SECRET")) // get jwt secret signature

	http.Handle("/", ValidateJWT(Home))
	http.HandleFunc("/jwt", GetJWT)
	fmt.Println("Server started at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
