package auth

import (
	"go-jwt/config"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var app *config.Config

// Newhelpers sets up appconfig for helpers
func NewAuth(a *config.Config) {
	app = a
}

type Claims struct {
	UserID   string `json:"userid"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

func CreateToken(userid, username, email, role string, expiration time.Duration) (string, error) {
	claims := &Claims{
		UserID:   userid,
		Username: username,
		Email:    email,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiration)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(app.Auth.SecretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ValidateToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(app.Auth.SecretKey), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}
