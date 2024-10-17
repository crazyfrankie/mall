package auth

import (
	"github.com/gin-gonic/contrib/cors"
	"github.com/gin-gonic/gin"
	"time"
)

func CORS() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
		// 是否允许带 cookie 之类的
		AllowCredentials: true,
		AllowedOrigins:   []string{"http://localhost:8080"},
		// 不加这一行 前端拿不到 token
		ExposedHeaders: []string{"x-jwt-token"},
		MaxAge:         12 * time.Hour,
	})
}
