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
		SecretKey       string // key used for JWT
		RefreshTokenExp time.Duration
		AccessTokenExp  time.Duration
		AesSecretKey    []byte //  key used for AES-GCM encryption
	}
}
