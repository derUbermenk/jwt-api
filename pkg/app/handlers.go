package app

import (
	"jwt-auth-gin/pkg/api"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

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
	}
}

func (s *Server) ValidateToken() gin.HandlerFunc {
	return func(c *gin.Context) {
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
