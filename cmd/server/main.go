package main

import (
	"fmt"
	"jwt-auth-gin/pkg/api"
	"jwt-auth-gin/pkg/app"
	"jwt-auth-gin/pkg/repository"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var Users = map[string]string{
	"user1": "password1",
	"user2": "password2",
}

var jwtKey = []byte("my_secret_key")

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "startup error encountered: %s \\n", err)
		os.Exit(1)
	}
}

func run() error {
	storage := repository.NewStorage(Users)

	// setup services
	authService := api.NewAuthService(storage) // will handle user authorization
	userService := api.NewUserService(storage) // will handle user manipulation

	// user gin router with logger and recovery middleware
	router := gin.Default()
	// also add global support for cors
	router.Use(cors.Default())

	server := app.NewServer(router, userService, authService)
	err := server.Run()

	if err != nil {
		return err
	}

	return nil
}
