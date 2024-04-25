package main

import (
	"fmt"
	"net/http"
)

// Handler is a struct that holds shared data and methods for handling requests.
type Handler struct {
	data string
}

// LoginHandler handles login requests.
func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	// Simulate populating data during login.
	h.data = "User123"
	fmt.Fprintf(w, "Login successful\n")
}

// ProfileHandler handles profile requests.
func (h *Handler) ProfileHandler(w http.ResponseWriter, r *http.Request) {
	// Access the shared data.
	fmt.Fprintf(w, "Profile data: %s\n", h.data)
}

// LogoutHandler handles logout requests.
func (h *Handler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// Clear the shared data.
	h.data = ""
	fmt.Fprintf(w, "Logout successful\n")
}

func main() {
	handler := &Handler{}

	// Register handlers.
	http.HandleFunc("/login", handler.LoginHandler)
	http.HandleFunc("/profile", handler.ProfileHandler)
	http.HandleFunc("/logout", handler.LogoutHandler)

	// Start the server.
	fmt.Println("Server listening on port 8080...")
	http.ListenAndServe(":8080", nil)
}
