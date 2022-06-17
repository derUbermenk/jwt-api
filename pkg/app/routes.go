package app

import "github.com/gin-gonic/gin"

func (s *Server) Routes() *gin.Engine {
	router := s.router

	router.POST("/signIn", s.SignIn())
	router.POST("/session/refresh", s.ValidateRefreshToken(), s.RefreshAccessToken())

	private := router.Group("private/")
	private.Use(s.ValidateUser())
	{
		private.GET("/user/greet/:username", s.GetUser())
	}

	return router
}
