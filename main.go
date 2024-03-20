package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

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
	r.Get("/refresh-token", refreshTokenHandler) // get new access token using the refresh token
	r.Get("/logout", Logout)
	r.Post("/generate-tokens", generateTokensHandler) // get 2 tokens: access + refresh
	// protected routes
	r.Route("/protected", func(mux chi.Router) {
		mux.Use(Auth)
		mux.Get("/user", func(w http.ResponseWriter, r *http.Request) { // protected/user
			secKey, err := getSecKey()
			if err != nil {
				log.Fatal(err)
			}
			value, err := ReadEncrypted(r, "refresh_token", secKey)
			if err != nil {
				ErrorJSON(w, err, http.StatusBadRequest)
				return
			}
			//// Validate the access token
			validToken, err := validateToken(value)
			if err != nil || !validToken.Valid {
				redirectURL := fmt.Sprintf("/refresh-token?redirect=%s", url.QueryEscape(r.URL.Path))
				http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
				return
			}
			//// Access claims from the token
			claims, ok := validToken.Claims.(*Claims)
			if !ok {
				ErrorJSON(w, err, http.StatusUnauthorized)
				return
			}
			//// send all tokens in a json format to the UI
			payload := JsonResponse{
				Error:   false,
				Message: fmt.Sprintf("Loggedin user of email %s and ID: %s", claims.Email, claims.UserID),
			}
			WriteJSON(w, http.StatusAccepted, payload)
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
