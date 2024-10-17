package auth

import (
	"github.com/gin-gonic/contrib/cors"
	"github.com/gin-gonic/gin"
	"time"
)

func CORS() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		AllowedOrigins:   []string{"http://localhost:8080"},
		ExposedHeaders:   []string{"x-jwt-token"},
		MaxAge:           12 * time.Hour,
	})
}
