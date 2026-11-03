package main

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TestJWTAuthenticator struct{}

const secret = "test"

var testClaims = jwt.MapClaims{
	"aud": "test-aud",
	"iss": "test-iss",
	"sub": 4,
	"exp": time.Now().Add(time.Hour).Unix(),
}

func NewTestJWTAuthenticator() *TestJWTAuthenticator {
	return &TestJWTAuthenticator{}
}

func (auth *TestJWTAuthenticator) GenerateToken(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, testClaims)
	return token.SignedString([]byte(secret))
}

func (a *TestJWTAuthenticator) ValidateToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(t *jwt.Token) (any, error) {
		return []byte(secret), nil
	})
}
