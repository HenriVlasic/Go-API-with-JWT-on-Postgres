package middleware

import (
	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(context *gin.Context) {
		tokenCookie, er := context.Request.Cookie("jwt")
		if er != nil {
			context.AbortWithStatusJSON(http.StatusUnauthorized, er)
			return
		}

		_, err := jwt.ParseWithClaims(tokenCookie.Value, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("SECRET")), nil
		})

		if err != nil {
			context.JSON((http.StatusUnauthorized), gin.H{
				"status_code":	http.StatusUnauthorized,
				"message":		"Unauthenticated",
				"error":		err.Error(),
			})
			return
		}
	}
}