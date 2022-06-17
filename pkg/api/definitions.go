package api

import (
	"github.com/golang-jwt/jwt/v4"
)

var jwtKey = []byte("my_secret_key")

type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

type AccessTokenClaims struct {
	TokenType string `json:"token_type"`
	Username  string `json:"username"`

	jwt.RegisteredClaims
}

type RefreshTokenClaims struct {
	TokenType string `json:"token_type"`
	Username  string `json:"username"`
	CustomKey string `json:"custom_key"`

	// always add this so that this custom claim can be considered of type jwt.Claims
	jwt.RegisteredClaims
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var (
	SignatureInvalidError   = jwt.ErrSignatureInvalid
	ExpiredAccessTokenError = jwt.ErrTokenExpired
)
