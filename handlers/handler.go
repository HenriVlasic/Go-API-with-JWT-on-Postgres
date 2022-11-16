package handler

import (
	"cyberzell.com/seguros/handlers/middleware"
	"cyberzell.com/seguros/models"
	"cyberzell.com/seguros/models/apperrors"

	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	userService		models.IUserService
	MaxBodyBytes	int64
}

// Config will hold services that will eventually be injected into this
// handler layer on handler intialization
type Config struct {
	R				*gin.Engine
	UserService		models.IUserService
	MaxBodyBytes	int64
}

// NewHandler initializes the handler with required injected services along with http routes
// Does not return as it deals directly with a reference to the gin Engine
func NewHandler(c *Config) {
	// Create a handler (which will later have injected service)
	h := &Handler {
		userService:	c.UserService,
		MaxBodyBytes:	c.MaxBodyBytes,
	}

	c.R.NoRoute(func(c *gin.Context){
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Page not found",
		})
	})

	// Create an account group
	userGroup := c.R.Group("api/user")

	userGroup.POST("/register", h.Register)
	userGroup.POST("/login", h.Login)
	userGroup.POST("/logout", h.Logout)

	authRoutes := c.R.Group("api/users").Use(middleware.Auth())
	authRoutes.GET("/:id", h.GetUserById)
}

func toFieldErrorResponse(c *gin.Context, field, message string) {
	c.JSON(http.StatusBadRequest, gin.H{
		"errors": []apperrors.FieldError{
			{
				Field:		field,
				Message:	message,
			},
		},
	})
}