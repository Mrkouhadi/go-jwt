package main

import (
	"go-jwt/internal/config"
	"go-jwt/internal/handlers"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func Routes(app *config.Config) http.Handler {

	r := chi.NewRouter()
	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	// Basic CORS
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins: []string{"http://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*", "http://localhost:3000"},
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
	// authentication endpoints
	r.Get("/refresh-token", handlers.Repo.RefreshTokenHandler) // get new access token using the refresh token
	r.Get("/logout", handlers.Repo.Logout)
	r.Post("/login", handlers.Repo.GenerateTokensHandler) // get 2 tokens: access + refresh
	// protected routes
	r.Route("/protected", func(mux chi.Router) {
		mux.Use(Auth)
		mux.Get("/profile", handlers.Repo.Profile)
	})
	return r
}
