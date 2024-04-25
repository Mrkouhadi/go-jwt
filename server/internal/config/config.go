package config

import "time"

type Config struct {
	Auth struct {
		UserData struct {
			UserID   string `json:"userid"`
			Username string `json:"username"`
			Email    string `json:"email"`
			Role     string `json:"role"`
		}
		SecretKey       string
		RefreshTokenExp time.Duration
		AccessTokenExp  time.Duration
	}
	PortNumber string
}
