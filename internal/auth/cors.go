package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		if origin == "http://localhost:8080" {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Exposed-Headers", "x-jwt-token")
			c.Header("Access-Control-Max-Age", "43200") // 12小时
		}

		if c.Request.Method == "OPTIONS" {
			c.Status(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
