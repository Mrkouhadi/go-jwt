package handlers

import "go-jwt/config"

// the repository used by the handlers
var Repo *Repository

// repository type
type Repository struct {
	App *config.Config
}

// NewRepo creates the new repository
func NewRepo(a *config.Config) *Repository {
	return &Repository{
		App: a,
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}
