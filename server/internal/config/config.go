package config

import (
	"log"
	"time"
)

type Config struct {
	PortNumber string
	InfoLog    *log.Logger
	ErrorLog   *log.Logger
	Auth       struct {
		UserData struct {
			UserID   string `json:"userid"`
			Username string `json:"username"`
			Email    string `json:"email"`
			Role     string `json:"role"`
		}
		SecretKey       string // key used for JWT
		RefreshTokenExp time.Duration
		AccessTokenExp  time.Duration
		AesSecretKey    []byte //  key used for AES-GCM encryption
	}
}
