package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/himanshuk42/gin-authentication/responses"
	"github.com/himanshuk42/gin-authentication/utils/token"
)

func JwtAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := token.TokenValid(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, responses.UserResponse{Status: http.StatusUnauthorized, Message: "unauthorized", Data: map[string]interface{}{"data": "User Unauthorized"}})
			return
		}
	}
}
