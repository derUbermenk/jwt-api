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
	AccessToken string `json:"access_token"`
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

		access_token, err := s.authService.MakeToken(cred)

		if err != nil {
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
				Data:    &AuthResponse{AccessToken: access_token},
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

		tkn_valid, current_user, err := s.authService.ValidateToken(token_string)

		if err != nil {
			c.JSON(
				http.StatusBadRequest,
				&GenericResponse{Status: false, Message: "bad request"},
			)
			return
		}

		if !tkn_valid {
			c.JSON(
				http.StatusUnauthorized,
				&GenericResponse{Status: false, Message: "User not authenticated"},
			)
			return
		}

		c.Set("current_user", current_user)
		c.Next()
	}
}

func (s *Server) GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		// get user given the usersname
		userName := c.Param("username")
		user, err := api.UserService.GetUser(userName)

		if err != nil {
			log.Printf("handler error: %v", err)
			c.JSON(http.StatusBadRequest, GenericResponse{Status: false, Message: err.Error()})
		}

		c.JSON(http.StatusOK, GenericResponse{Status: true, Message: "User retrieved", Data: user})
	}
}
