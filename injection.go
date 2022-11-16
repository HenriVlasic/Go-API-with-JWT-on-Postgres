package main

import (
	handler "cyberzell.com/seguros/handlers"
	repository "cyberzell.com/seguros/repositories"
	"cyberzell.com/seguros/services"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func inject(d *dataSources) (*gin.Engine, error) {
	log.Println("Injecting data sources")

	/*
	* repository layer
	*/
	userRepository := repository.NewUserRepository(d.DB)

	origin := os.Getenv("CORS_ORIGIN")

	/*
	* service layer
	*/
	userService := services.NewUserService(&services.USConfig{
		UserRepository: userRepository,
	})

	router := gin.Default()

	c := cors.New(cors.Config{
		AllowOrigins: 		[]string{origin},
		AllowCredentials:	true,
		AllowMethods:		[]string{"GET", "POST", "PUT", "DELETE"},
	})

	router.Use(c)

	maxBodyBytes := os.Getenv("MAX_BODY_BYTES")
	mbb, err := strconv.ParseInt(maxBodyBytes, 0, 64)
	if err != nil {
		return nil, fmt.Errorf("could not parse MAX_BODY_BYTES as int: %w", err)
	}

	handler.NewHandler(&handler.Config{
		R:				router,
		UserService: 	userService,
		MaxBodyBytes:	mbb,
	})

	return router, nil
}