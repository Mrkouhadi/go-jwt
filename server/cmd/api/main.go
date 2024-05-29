package main

import (
	"encoding/hex"
	"fmt"
	"go-jwt/internal/auth"
	"go-jwt/internal/config"
	"go-jwt/internal/handlers"
	"go-jwt/internal/helpers"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

var app = config.Config{}
var infoLog *log.Logger
var errorLog *log.Logger

func main() {
	// port number of our server
	app.PortNumber = ":8080"

	// set up loggers
	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime) // \t means tab (bunch of spaces)
	app.InfoLog = infoLog
	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	// jwt configs
	scretKey := goDotEnvVariable("SECRET_JWT_KEY")
	aesScretKey := goDotEnvVariable("AES_SECRECT_KEY")
	app.Auth.SecretKey = scretKey                 // jwt very secret
	app.Auth.AccessTokenExp = time.Minute * 30    // 30 minutes
	app.Auth.RefreshTokenExp = time.Hour * 24 * 7 // 7 days

	// settin gup a secrect key for AES-GCM encryption: because we encrypt jwt tokens before storing them in cookies
	secKey, err := hex.DecodeString(aesScretKey)
	if err != nil {
		log.Fatal(err)
	}
	app.Auth.AesSecretKey = secKey // key used for AES-GCM encryption

	// set up handlers
	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	// set up auth
	auth.NewAuth(&app)

	// helpers
	helpers.Newhelpers(&app)

	// set up the server
	srv := &http.Server{
		Addr:              app.PortNumber,
		Handler:           Routes(&app),
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		ReadHeaderTimeout: 15 * time.Second,
	}
	fmt.Println("Server running on ", app.PortNumber)
	log.Fatal(srv.ListenAndServe())
}

// goDotEnvVariable is a helper to get .env data
func goDotEnvVariable(key string) string {
	// load .env file
	// if we have specific file names like this: system.env & others.env we can just do : godotenv.Load("system.env")
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	return os.Getenv(key)
}
