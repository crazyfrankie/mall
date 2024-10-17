package auth

import (
	"github.com/gin-gonic/contrib/cors"
	"github.com/gin-gonic/gin"
	"time"
)

//func CORS() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		origin := c.Request.Header.Get("Origin")
//		if origin == "http://localhost:8080" {
//			c.Header("Access-Control-Allow-Origin", origin)
//			c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
//			c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
//			c.Header("Access-Control-Allow-Credentials", "true")
//			c.Header("Exposed-Headers", "x-jwt-token")
//			c.Header("Access-Control-Max-Age", "43200") // 12小时
//		}
//
//		// 如果是预检请求，返回204 No Content
//		if c.Request.Method == "OPTIONS" {
//			c.Status(http.StatusNoContent)
//			return
//		}
//
//		c.Next()
//	}
//}

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
