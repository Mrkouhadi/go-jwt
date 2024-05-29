package helpers

import (
	"encoding/json"
	"errors"
	"go-jwt/internal/config"
	"html"
	"io"
	"net/http"
	"regexp"
	"strings"
)

var app *config.Config

// Newhelpers sets up appconfig for helpers
func Newhelpers(a *config.Config) {
	app = a
}

// validateAndSanitizeEmail validates and sanitizes the email input
func ValidateAndSanitizeEmail(email string) (string, error) {
	// Trim leading and trailing whitespace
	email = strings.TrimSpace(email)
	// Escape HTML special characters
	email = html.EscapeString(email)
	// Regular expression for email validation
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	match, _ := regexp.MatchString(emailRegex, email)
	if !match {
		return "", errors.New("invalid email format")
	}
	// Regular expression for malicious code detection
	maliciousRegex := `<script|javascript:|onerror=|onload=|alert\(|confirm\(|<\/script>`
	if matched, _ := regexp.MatchString(maliciousRegex, email); matched {
		return "", errors.New("potentially malicious input detected")
	}
	return email, nil
}

// validateAndSanitizePassword validates and sanitizes the password input
func ValidateAndSanitizePassword(password string) (string, error) {
	// Trim leading and trailing whitespace
	password = strings.TrimSpace(password)

	// Escape HTML special characters
	password = html.EscapeString(password)

	// Check password length
	if len(password) < 8 {
		return "", errors.New("password must be at least 8 characters long")
	}

	// Regular expression for malicious code detection
	maliciousRegex := `<script|javascript:|onerror=|onload=|alert\(|confirm\(|<\/script>`
	if matched, _ := regexp.MatchString(maliciousRegex, password); matched {
		return "", errors.New("potentially malicious input detected")
	}
	// Regular expressions for required character types
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	hasSpecial := regexp.MustCompile(`[!@#$%^&*()_+=\-/?~}{\[\];:.,'']`).MatchString(password)

	// Check for required character types
	if !hasLower || !hasUpper || !hasNumber || !hasSpecial {
		return "", errors.New("password must contain at least 1 uppercase letter, 1 lowercase letter, 1 number, and 1 special character")
	}

	return password, nil
}

// /////////////////////////////////////// JSON /////////////////////////////////

type JsonResponse struct {
	Error   bool        `json:"error"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func ReadJSON(w http.ResponseWriter, r *http.Request, data interface{}) error {
	// limit the size of the uploaded json file
	maxByte := 1048576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxByte))
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(data)
	if err != nil {
		app.ErrorLog.Println("Error decoding body data")
		return err
	}
	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		app.ErrorLog.Println("body should have only one single JSON value")
		return errors.New("body should have only one single JSON value")
	}
	return nil
}

func WriteJSON(w http.ResponseWriter, status int, data interface{}, headers ...http.Header) error {
	out, err := json.Marshal(data)
	if err != nil {
		app.ErrorLog.Println("Error marshalling JSON data")
		return err
	}
	if len(headers) > 0 {
		for key, value := range headers[0] {
			w.Header()[key] = value
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(out)
	if err != nil {
		app.ErrorLog.Println("error writing JSON data")
		return err
	}
	return nil
}

func ErrorJSON(w http.ResponseWriter, err error, status ...int) error {
	statusCode := http.StatusBadRequest
	if len(status) > 0 {
		statusCode = status[0]
	}
	var payload JsonResponse
	payload.Error = true
	payload.Message = err.Error()
	return WriteJSON(w, statusCode, payload)
}

// /////////////////////////////////////// USAGE /////////////////////////////////
// example:
/*
	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	err := ReadJSON(w,r, &requestPayload)
	if err != nil {
		ErrorJSON(w, err, http.StatusBadRequest)
		return
	}
	payload := JsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Logged in user %s", user.Email),
		Data:    user,
	}
	WriteJSON(w, http.StatusAccepted, payload)
}
*/
