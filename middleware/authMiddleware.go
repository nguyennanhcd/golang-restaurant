package middlewares

import (
	"golang-restaurant-management/helper"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {

		clientToken := strings.TrimSpace(c.Request.Header.Get("Authorization"))
		if clientToken == "" {
			c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "No token provided"})
			c.Abort()
			return
		}

		if strings.HasPrefix(clientToken, "Bearer ") {
			clientToken = strings.TrimPrefix(clientToken, "Bearer ")
			clientToken = strings.TrimSpace(clientToken)
		} else {
			c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Invalid token format, expected 'Bearer <token>'"})
			c.Abort()
			return
		}

		claims, errMsg := helper.ValidateToken(clientToken)
		if errMsg != "" {
			c.JSON(http.StatusUnauthorized, ErrorResponse{Error: errMsg})
			c.Abort()
			return
		}

		if claims.TokenType != "access" {
			c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Invalid token type, access token required"})
			c.Abort()
			return
		}

		c.Set("email", claims.Email)
		c.Set("first_name", claims.First_name)
		c.Set("last_name", claims.Last_name)
		c.Set("uid", claims.Uid)

		c.Next()
	}
}
