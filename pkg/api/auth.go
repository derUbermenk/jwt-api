package api

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"log"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type AuthenticationService interface {
	ValidateAccessToken(access_token string) (isValid bool, username string, err error)
	GenerateAccessToken(creds Credentials) (access_token string, err error)
	ValidateRefreshToken(refresh_token string) (isValid bool, username string, err error)
	GenerateRefreshToken(creds Credentials) (refresh_token string, err error)
	ValidateCredentials(creds Credentials) (valid bool, err error)
}

type AuthRepository interface {
	GetUser(userName string) (user User, err error)
}

type authService struct {
	jwtKey  []byte
	storage AuthRepository
}

func NewAuthService(jwtKey []byte, authRepo AuthRepository) AuthenticationService {
	return &authService{
		jwtKey:  jwtKey,
		storage: authRepo,
	}
}

func (a *authService) ValidateAccessToken(token_string string) (validity bool, username string, err error) {
	claims := &AccessTokenClaims{}

	// Parse the JWT string and store the result in `claims`.
	// Note that we are passing the key in this method as well. This method will return an error
	// if the token is invalid (if it has expired according to the expiry time we set on sign in),
	// or if the signature does not match
	tkn, err := jwt.ParseWithClaims(token_string, claims, func(token *jwt.Token) (interface{}, error) { return a.jwtKey, nil })

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return false, "", nil
		}

		// err could also be caused when the jwt token is expired
		// but compare the error with jwt.ErrTokenExpired does not work
		// to augment this, I resolved to check the error value instead
		if strings.Contains(err.Error(), jwt.ErrTokenExpired.Error()) {
			return false, "", nil
		}

		return false, "", err
	}

	if !tkn.Valid {
		return false, "", nil
	}

	return true, claims.Username, nil
}

func (a *authService) GenerateAccessToken(creds Credentials) (access_token_string string, err error) {
	expirationTime := time.Now().Add(5 * time.Minute)

	claims := &AccessTokenClaims{
		TokenType: "access",
		Username:  creds.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	access_token_string, err = token.SignedString(a.jwtKey)

	if err != nil {
		log.Printf("Service error: %v", err)
		return access_token_string, err
	}

	return access_token_string, nil
}

func (a *authService) ValidateRefreshToken(token_string string) (validity bool, username string, err error) {
	claims := &RefreshTokenClaims{}

	tkn, err := jwt.ParseWithClaims(token_string, claims, func(token *jwt.Token) (interface{}, error) { return a.jwtKey, nil })

	if err != nil {
		log.Printf("service errror: %v", err)
		if err == jwt.ErrSignatureInvalid {
			// disregard the error, this just means that
			// the user is not authorized, maybe the
			// token was tampered
			return false, "", nil
		}

		return false, "", err
	}

	if !tkn.Valid || claims.Username == "" {
		return false, "", nil
	}

	// get the user
	user, err := a.storage.GetUser(claims.Username)

	if err != nil {
		log.Printf("service errror: %v", err)
		return false, "", err
	}

	if claims.CustomKey != a.GenerateCustomKey(user.Username, user.Password) {
		return false, "", nil
	}

	return true, claims.Username, nil
}

func (a *authService) GenerateRefreshToken(creds Credentials) (refresh_token_string string, err error) {
	cusKey := a.GenerateCustomKey(creds.Username, creds.Password)

	claims := &RefreshTokenClaims{
		TokenType: "refresh",
		Username:  creds.Username,
		CustomKey: cusKey,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	refresh_token_string, err = token.SignedString(a.jwtKey)

	if err != nil {
		log.Printf("Service error: %v", err)
		return "", err
	}

	return refresh_token_string, nil
}

func (a *authService) ValidateCredentials(creds Credentials) (valid bool, err error) {
	user, err := a.storage.GetUser(creds.Username)

	if err != nil {
		log.Printf("Service error: %v", err)
		return false, err
	}

	if user.Password != creds.Password {
		return false, nil
	}

	return true, nil
}

// generates some hashed key from the
func (a *authService) GenerateCustomKey(username, password string) (customKey string) {

	// create a has object h
	// 	https://betterprogramming.pub/a-short-guide-to-hashing-in-go-e8bb0173e97e
	// h.Write add more data, in our case the byte representation of the username to the hash
	// encode the byte representation of the current hash to a string
	// to do this we h.Sum(nil) essentially adding nothing to the hash byte representation

	h := hmac.New(sha256.New, []byte(password))
	h.Write([]byte(username))

	customKey = hex.EncodeToString(h.Sum(nil))
	return
}
