package app

import "github.com/gin-gonic/gin"

func (s *Server) Routes() *gin.Engine {
	router := s.router

	public := router.Group("public/")
	{
		public.GET("/user/greet/:username", s.userService.PublicGreet())
	}

	private := router.Group("private/")
	{
		private.GET("/user/greet/:username", s.userService.PrivateGreet())
	}

	return router
}
