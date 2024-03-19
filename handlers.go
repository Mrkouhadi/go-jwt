package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

// generateTokensHandler represents LOGIN handler in real world scenario
func generateTokensHandler(w http.ResponseWriter, r *http.Request) {
	//// Extract username from the request, validate credentials, etc. ////
	var Person struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	err := ReadJSON(w, r, &Person)
	if err != nil {
		ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}
	//// compare credentials with Main database data.
	//
	//
	//
	//// if the credentials match what we have in the main database we issue the Tokens
	//// Create an access token
	accessToken, err := createToken(Person.Password, Person.Email, accessTokenExp)
	if err != nil {
		ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	//// Create a refresh token
	refreshToken, err := createToken(Person.Password, Person.Email, refreshTokenExp)
	if err != nil {
		ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	//// storing the tokens in cookies
	// Create multiple cookies
	cookies := []http.Cookie{
		{
			Name:     "access_token",
			Value:    accessToken,
			HttpOnly: true,
			// Secure:   true,
			Path:     "/",
			MaxAge:   3600,
			SameSite: http.SameSiteLaxMode,
		},
		{
			Name:     "refresh_token",
			Value:    refreshToken,
			HttpOnly: true,
			// Secure:   true,
			Path:     "/",
			MaxAge:   3600,
			SameSite: http.SameSiteLaxMode,
		},
	}
	// get the secret key
	secKey, err := getSecKey()
	if err != nil {
		ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}
	// Set each cookie
	for _, cookie := range cookies {
		err := WriteEncrypted(w, cookie, secKey)
		if err != nil {
			ErrorJSON(w, err, http.StatusInternalServerError)
			return
		}
	}
	//// send all tokens in a json format to the UI
	payload := JsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Logged in user %s", Person.Email),
	}
	WriteJSON(w, http.StatusAccepted, payload)
}

func refreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	secKey, err := getSecKey()
	if err != nil {
		log.Fatal(err)
	}
	value, err := ReadEncrypted(r, "refresh_token", secKey)
	if err != nil {
		ErrorJSON(w, err, http.StatusBadRequest)
		return
	}
	// Validate the refresh token
	validToken, err := validateToken(value)
	if err != nil {
		ErrorJSON(w, err, http.StatusUnauthorized)
		return
	}
	// Access claims from the token
	claims, ok := validToken.Claims.(*Claims)
	if !ok {
		ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}
	// Check if the refresh token is expired
	expiresAtUnix := claims.ExpiresAt.Time.Unix()
	if time.Now().Unix() >= expiresAtUnix {
		ErrorJSON(w, err, http.StatusUnauthorized)
		return
	}
	// Create a new access token
	newAccessToken, err := createToken(claims.Password, claims.Email, accessTokenExp)
	if err != nil {
		ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	// UPDATE the access token in the cookie
	NewAccTokenCookie := http.Cookie{
		Name:     "access_token",
		Value:    newAccessToken,
		HttpOnly: true,
		// Secure:   true,
		Path:     "/",
		MaxAge:   3600,
		SameSite: http.SameSiteLaxMode,
	}
	err = WriteEncrypted(w, NewAccTokenCookie, secKey)
	if err != nil {
		ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}
	//// REDIRECT THE USER BACK TO THE ORIGINAL LINK
	redirectURL := r.URL.Query().Get("redirect")
	if redirectURL == "" {
		// If no redirect URL is provided, redirect to a default URL
		redirectURL = "/"
	}
	http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
}

// LOGOUT: handler for logging out by deleting all tokesn from cookies, and remove redis credentials
func Logout(w http.ResponseWriter, r *http.Request) {
	// Create a new cookie with the same name as the one you want to delete
	// and set its expiration time to a time in the past (e.g., one hour ago).
	deletedCookies := []http.Cookie{
		{
			Name:     "refresh_token",
			Value:    "", // Setting the value to empty is optional but recommended
			Expires:  time.Now().Add(-time.Hour),
			HttpOnly: true, // Set other cookie attributes as needed
			Path:     "/",  // Set the cookie's path
		},
		{
			Name:     "access_token",
			Value:    "", // Setting the value to empty is optional but recommended
			Expires:  time.Now().Add(-time.Hour),
			HttpOnly: true, // Set other cookie attributes as needed
			Path:     "/",  // Set the cookie's path
		},
	}
	// Set the cookie on the response with an expiration time in the past
	for _, deletedCookie := range deletedCookies {
		http.SetCookie(w, &deletedCookie)
	}
	// redirect the user
	http.Redirect(w, r, "/", http.StatusFound)
}
