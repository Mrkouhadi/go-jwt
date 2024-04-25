package main

import (
	"fmt"
	"go-jwt/auth"
	"go-jwt/config"
	"go-jwt/handlers"
	"log"
	"net/http"
	"time"
)

func main() {
	app := config.Config{}
	app.PortNumber = ":8080"
	app.Auth.SecretKey = "Fd~+?/-S@c?re$t99~~!!~~8R-Fr5?>?" // very secret; u should use .env
	app.Auth.AccessTokenExp = time.Minute * 30              // 30 minutes
	app.Auth.RefreshTokenExp = time.Hour * 24 * 7           // 7 days

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

// login link: http://localhost:8080/generate-tokens    POST
// {
//     "username":"bryan",
//     "email":"bryan@kouhadi.com"
// }
// refresh the access token when it's expired: http://localhost:8080/refresh-token  GET
// protected link:  http://localhost:8080/protected/user   GET
// logout link: http://localhost:8080/logout  GET
