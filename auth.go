package main

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const (
	secretKey       = "Fd~+?/-S@c?ret99~~!!~~8R-Fr5?>?" // very secret
	accessTokenExp  = time.Minute * 15                  //FIXME: make it 15 minutes
	refreshTokenExp = time.Hour * 24 * 7                // 7 days
)

type Claims struct {
	Email  string `json:"email"`
	UserID string `json:"userid"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func createToken(userid, email string, expiration time.Duration) (string, error) {
	claims := &Claims{
		Email:  email,
		UserID: userid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiration)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func validateToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}
