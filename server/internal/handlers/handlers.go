package handlers

import (
	"fmt"
	"go-jwt/internal/auth"
	"go-jwt/internal/helpers"
	"net/http"
	"time"
)

// Login represents LOGIN handler in real world scenario
func (m *Repository) Login(w http.ResponseWriter, r *http.Request) {
	//// Extract username from the request, validate credentials, etc.
	var Person struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	err := helpers.ReadJSON(w, r, &Person)
	if err != nil {
		helpers.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}
	validatedEmail, err := helpers.ValidateAndSanitizeEmail(Person.Email)
	if err != nil {
		helpers.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}
	validatedPassword, err := helpers.ValidateAndSanitizePassword(Person.Password)
	if err != nil {
		helpers.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	fmt.Println("Sanitized and Validated Password: ", validatedPassword)

	////
	//
	//  TODO: compare credentials against database data.
	//
	////

	// Create an access token
	userid := "id_ftdrs42671hdcn" // get it from the database
	username := "bryan kouhadi"
	email := validatedEmail
	role := "user"
	accessToken, err := auth.CreateToken(userid, username, email, role, m.App.Auth.AccessTokenExp)
	if err != nil {
		helpers.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	//// Create a refresh token
	refreshToken, err := auth.CreateToken(userid, username, email, role, m.App.Auth.RefreshTokenExp)
	if err != nil {
		helpers.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	//// storing the tokens in cookies
	// Create multiple cookies
	cookies := []http.Cookie{
		{
			Name:     "access_token",
			Value:    accessToken,
			HttpOnly: true,
			Secure:   true,
			Path:     "/",
			MaxAge:   3600,
			SameSite: http.SameSiteLaxMode,
		},
		{
			Name:     "refresh_token",
			Value:    refreshToken,
			HttpOnly: true,
			Secure:   true,
			Path:     "/",
			MaxAge:   3600,
			SameSite: http.SameSiteLaxMode,
		},
	}

	// Set each cookie
	for _, cookie := range cookies {
		err := auth.WriteEncrypted(w, cookie, m.App.Auth.AesSecretKey)
		if err != nil {
			helpers.ErrorJSON(w, err, http.StatusInternalServerError)
			return
		}
	}
	//// send all tokens in a json format to the UI
	payload := helpers.JsonResponse{
		Error:   false,
		Message: "Logged in successfully",
		Data:    struct{ userid, username string }{username, userid},
	}
	helpers.WriteJSON(w, http.StatusAccepted, payload)
}

// RefreshTokenHandler refreshes the raccess token
func (m *Repository) RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	value, err := auth.ReadEncrypted(r, "refresh_token", m.App.Auth.AesSecretKey)
	if err != nil {
		helpers.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}
	// Validate the refresh token
	validToken, err := auth.ValidateToken(value)
	if err != nil {
		helpers.ErrorJSON(w, err, http.StatusUnauthorized)
		return
	}
	// Access claims from the token
	claims, ok := validToken.Claims.(*auth.Claims)
	if !ok {
		helpers.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}
	// Check if the refresh token is expired
	expiresAtUnix := claims.ExpiresAt.Time.Unix()
	if time.Now().Unix() >= expiresAtUnix {
		helpers.ErrorJSON(w, err, http.StatusUnauthorized)
		return
	}
	// Create a new access token
	newAccessToken, err := auth.CreateToken(claims.UserID, claims.Username, claims.Email, claims.Role, m.App.Auth.AccessTokenExp)
	if err != nil {
		helpers.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	// UPDATE the access token in the cookie
	NewAccTokenCookie := http.Cookie{
		Name:     "access_token",
		Value:    newAccessToken,
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		MaxAge:   3600,
		SameSite: http.SameSiteLaxMode,
	}
	err = auth.WriteEncrypted(w, NewAccTokenCookie, m.App.Auth.AesSecretKey)
	if err != nil {
		helpers.ErrorJSON(w, err, http.StatusInternalServerError)
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
func (m *Repository) Logout(w http.ResponseWriter, r *http.Request) {
	// Create a new cookie with the same name as the one you want to delete
	// and set its expiration time to a time in the past (e.g., one hour ago).
	deletedCookies := []http.Cookie{
		{
			Name:     "refresh_token",
			Value:    "", // Setting the value to empty is optional but recommended
			Expires:  time.Unix(0, 0),
			HttpOnly: true, // Set other cookie attributes as needed
			Path:     "/",  // Set the cookie's path
		},
		{
			Name:     "access_token",
			Value:    "", // Setting the value to empty is optional but recommended
			Expires:  time.Unix(0, 0),
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

// Profile is a protected handler shoing the details of the logged user
func (m *Repository) Profile(w http.ResponseWriter, r *http.Request) { // protected/user
	//// send all tokens in a json format to the UI
	payload := helpers.JsonResponse{
		Error:   false,
		Message: "You are Allowed here",
	}
	helpers.WriteJSON(w, http.StatusAccepted, payload)
}
