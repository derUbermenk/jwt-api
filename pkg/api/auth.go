package api

import (
	"log"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type AuthenticationService interface {
	ValidateToken(token_string string) (isValid bool, username string, err error)
	MakeToken(creds Credentials) (tokenString string, err error)
	ValidateCredentials(creds Credentials) (valid bool, err error)
}

type AuthRepository interface {
	GetUser(userName string) (user User, err error)
}

type authService struct {
	jwtKey  []byte
	storage AuthRepository
}

func (a *authService) ValidateToken(token_string string) (validity bool, username string, err error) {
	claims := &Claims{}

	// Parse the JWT string and store the result in `claims`.
	// Note that we are passing the key in this method as well. This method will return an error
	// if the token is invalid (if it has expired according to the expiry time we set on sign in),
	// or if the signature does not match
	tkn, err := jwt.ParseWithClaims(token_string, claims, func(token *jwt.Token) (interface{}, error) { return a.jwtKey, nil })

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return false, "", nil
		}

		return false, "", err
	}

	if !tkn.Valid {
		return false, "", nil
	}

	return true, claims.Username, nil
}

func (a *authService) MakeToken(creds Credentials) (tokenString string, err error) {
	expirationTime := time.Now().Add(5 * time.Minute)

	claims := &Claims{
		Username: creds.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err = token.SignedString(a.jwtKey)

	if err != nil {
		log.Printf("Service error: %v", err)
		return tokenString, err
	}

	return tokenString, nil
}

func (a *authService) ValidateCredentials(creds Credentials) (valid bool, err error) {
	user, err := a.storage.GetUser(creds.Username)

	if err != nil {
		log.Printf("Service error: %v", err)
		return false, err
	}

	if (user == User{}) {
		return false, nil
	} else if user.Password != creds.Password {
		return false, nil
	}

	return true, nil
}
