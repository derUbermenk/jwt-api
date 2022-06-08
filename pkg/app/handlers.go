package app

import (
	"jwt-auth-gin/pkg/api"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GenericResponse struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (s *Server) SignIn() gin.HandlerFunc {
	return func(c *gin.Context) {
		// decode credentials
		// credentials {username, password}
		// get the user with username from database
		// if user !exist ;
		// Generic Response
		//	return
		// compare user password with credentials password
		// if user password and credentials match
		// sign in user
		//
		// create a jwt token
		// store it as cookie
		// return
		var cred api.Credentials

		err := c.ShouldBindJSON(&cred)

		if err != nil {
			c.JSON(
				http.StatusBadRequest,
				&GenericResponse{
					Status:  false,
					Message: "Bad Request",
				},
			)

			return
		}

		cred_valid, err := s.authService.ValidateCredentials(cred)

		if err != nil {
			c.JSON(
				http.StatusBadRequest,
				&GenericResponse{
					Status:  false,
					Message: "Bad Request",
				},
			)

			return
		}

		if !cred_valid {
			c.JSON(
				http.StatusUnauthorized,
				&GenericResponse{
					Status:  false,
					Message: "Invalid Credentials",
				},
			)
		}

		access_token, err := s.authService.GenerateAccessToken(cred)

		if err != nil {
			log.Printf("%v", err)

			c.JSON(
				http.StatusInternalServerError,
				&GenericResponse{
					Status: false, Message: "Internal Server Error",
				},
			)

			return
		}

		refresh_token, err := s.authService.GenerateRefreshToken(cred)

		if err != nil {
			log.Printf("%v", err)
			c.JSON(
				http.StatusInternalServerError,
				&GenericResponse{
					Status: false, Message: "Internal Server Error",
				},
			)

			return
		}

		c.JSON(
			http.StatusOK,
			&GenericResponse{
				Status:  true,
				Message: "Signed in successfully",
				Data:    &AuthResponse{AccessToken: access_token, RefreshToken: refresh_token},
			},
		)
	}
}

func (s *Server) ValidateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		token_string := c.GetHeader("AccessToken")

		// token_string refers to the base encoded jwt
		if token_string == "" {
			c.JSON(
				http.StatusBadRequest,
				&GenericResponse{Status: false, Message: "token not present"},
			)
			return
		}

		tkn_valid, current_user, err := s.authService.ValidateAccessToken(token_string)

		if err != nil {
			log.Printf("Internal Server Error: %v ", err)

			if err == api.ExpiredAccessTokenError {
				c.JSON(
					http.StatusInternalServerError,
					&GenericResponse{Status: false, Message: "expired access token"},
				)
			} else {
				c.JSON(
					http.StatusInternalServerError,
					&GenericResponse{Status: false, Message: "bad request"},
				)
			}

			c.Abort()
			return
		}

		if !tkn_valid {
			c.JSON(
				http.StatusUnauthorized,
				&GenericResponse{Status: false, Message: "User not authenticated"},
			)
			c.Abort()
			return
		}

		c.Set("current_user", current_user)
		c.Next()
	}
}

func (s *Server) RefreshAccessToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		// assuming the refresh token is confirmed to be valid
		// 		generate an access token

		// decode refresh token
		// get the user details
		// user details could be from refresh token claims
		// use the user.Username and user.Password as Credential variables
		// pass credentials into generate access token
		// return that access token

		var refresh_token_string string
		var creds api.Credentials

		access_token_string, err := s.authService.GenerateAccessToken(creds)

		if err != nil {
			log.Printf("%v", err)

			c.JSON(
				http.StatusInternalServerError,
				&GenericResponse{
					Status:  false,
					Message: "Internal Server Error",
				},
			)

			return
		}

		c.JSON(
			http.StatusOK,
			&GenericResponse{
				Status:  true,
				Message: "Access Token Refreshed",
				Data:    &AuthResponse{AccessToken: access_token_string, RefreshToken: refresh_token_string},
			},
		)
	}
}

func (s *Server) GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		// get user given the usersname
		userName := c.Param("username")
		user, err := s.userService.GetUser(userName)

		if err != nil {
			log.Printf("handler error: %v", err)
			c.JSON(http.StatusBadRequest, GenericResponse{Status: false, Message: err.Error()})
		}

		c.JSON(http.StatusOK, GenericResponse{Status: true, Message: "User retrieved", Data: user})
	}
}
