package app

import (
	"jwt-auth-gin/pkg/api"
	"log"

	"github.com/gin-gonic/gin"
)

type Server struct {
	router      *gin.Engine
	userService api.UserService
	authService api.AuthenticationService
}

func NewServer(router *gin.Engine, userService api.UserService, authService api.AuthenticationService) *Server {
	return &Server{
		router:      router,
		userService: userService,
		authService: authService,
	}
}

func (s *Server) Run() error {
	// initiate routes
	r := s.Routes()

	// run the server through the router
	err := r.Run()

	if err != nil {
		log.Printf("Server - there was an error calling Run on router: %v", err)
		return err
	}

	return nil
}
