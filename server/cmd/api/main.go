package main

import (
	"fmt"
	"go-jwt/internal/auth"
	"go-jwt/internal/config"
	"go-jwt/internal/handlers"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	app := config.Config{}
	scretKey := goDotEnvVariable("SECRET_JWT_KEY")
	app.Auth.SecretKey = scretKey // very secret

	app.PortNumber = ":8080"
	app.Auth.AccessTokenExp = time.Minute * 30    // 30 minutes
	app.Auth.RefreshTokenExp = time.Hour * 24 * 7 // 7 days

	// set up handlers
	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)
	// set up auth
	auth.NewAuth(&app)
	// set up the server
	srv := &http.Server{
		Addr:    app.PortNumber,
		Handler: Routes(&app),
	}
	fmt.Println("Server running on :8080")
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
