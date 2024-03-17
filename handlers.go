package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func generateTokensHandler(w http.ResponseWriter, r *http.Request) {
	// You may want to store the refresh token securely on the server for later use
	w.Header().Set("Content-Type", "application/json")
	//// Extract username from the request, validate credentials, etc. ////
	var Person struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	err := json.NewDecoder(r.Body).Decode(&Person)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Create access token
	accessToken, err := createToken(Person.Password, Person.Email, accessTokenExp)
	if err != nil {
		http.Error(w, "Error creating access token", http.StatusInternalServerError)
		return
	}

	// Create refresh token
	refreshToken, err := createToken(Person.Password, Person.Email, refreshTokenExp)
	if err != nil {
		http.Error(w, "Error creating refresh token", http.StatusInternalServerError)
		return
	}

	// Send tokens as JSON response
	response := map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}

	/// storing the tokens in cookies
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

	secKey, err := getSecKey()
	if err != nil {
		log.Fatal(err)
	}
	// Set each cookie
	for _, cookie := range cookies {
		// Use the WriteSigned() function, passing in the secret key as the final argument
		err := WriteEncrypted(w, cookie, secKey)
		if err != nil {
			http.Error(w, "server error", http.StatusInternalServerError)
			return
		}
	}
	// send all tokens in a json format to the UI
	json.NewEncoder(w).Encode(response)
}

func refreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// extracting the refresh token from the cookie(httpOnly cookie)
	secKey, err := getSecKey()
	if err != nil {
		log.Fatal(err)
	}
	value, err := ReadEncrypted(r, "refresh_token", secKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Validate the refresh token
	validToken, err := validateToken(value)
	if err != nil {
		http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
		return
	}
	// Access claims from the token
	claims, ok := validToken.Claims.(*Claims)
	if !ok {
		http.Error(w, "Error getting claims from token", http.StatusInternalServerError)
		return
	}
	// Check if the refresh token is expired
	if time.Until(time.Unix(claims.ExpiresAt, 0)) <= 0 {
		http.Error(w, "Refresh token has been expired", http.StatusUnauthorized)
		return
	}
	// Create a new access token
	newAccessToken, err := createToken(claims.Password, claims.Email, accessTokenExp)
	if err != nil {
		http.Error(w, "Error creating new access token", http.StatusInternalServerError)
		return
	}
	// Send the new access token as JSON response
	response := map[string]string{
		"access_token":  newAccessToken,
		"refresh_token": value,
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
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}
	// send json data to the UI
	json.NewEncoder(w).Encode(response)
}

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
	http.Redirect(w, r, "/", http.StatusFound)
}
