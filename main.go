package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

var client = NewRedisClient()

func main() {

	r := chi.NewRouter()
	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Basic CORS
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins: []string{"http://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) { // protected/user
		w.Write([]byte("HOME PAGE"))
	})
	r.Post("/generate-tokens", generateTokensHandler) // get 2 tokens: access + refresh
	r.Get("/refresh-token", refreshTokenHandler)      // get new access token using the refresh token
	r.Get("/logout", Logout)

	// protected routes
	r.Route("/protected", func(mux chi.Router) {
		mux.Use(Auth)
		mux.Get("/user", func(w http.ResponseWriter, r *http.Request) { // protected/user
			pass, err := ReadCredentials(client, "bryan@kouhadi.com")
			if err != nil {
				fmt.Println("error reading redis data in /protected/user: ", err)
			}
			fmt.Println("passs: ", pass)
			w.Write([]byte("User Profile-  password: " + pass))
		})
	})
	fmt.Println("Server running on :8080")
	http.ListenAndServe(":8080", r)
}

// login link: http://localhost:8080/generate-tokens    POST
// {
//     "username":"bryan",
//     "email":"bryan@kouhadi.com"
// }
// refresh the access token when it's expired: http://localhost:8080/refresh-token  GET
// protected link:  http://localhost:8080/protected/user   GET
// logout link: http://localhost:8080/logout  GET
