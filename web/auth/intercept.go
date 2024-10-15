package middleware

import (
	"github.com/gin-gonic/gin"

	"mall/middleware/jwt"
)

type LoginJWTMiddlewareBuilder struct {
	paths  map[string]struct{}
	jwtHdl *jwt.TokenHandler
}

func NewLoginJWTMiddlewareBuilder(jwtHdl *jwt.TokenHandler) *LoginJWTMiddlewareBuilder {
	return &LoginJWTMiddlewareBuilder{
		paths:  make(map[string]struct{}),
		jwtHdl: jwtHdl,
	}
}

func (b *LoginJWTMiddlewareBuilder) IgnorePath(path string) *LoginJWTMiddlewareBuilder {
	b.paths[path] = struct{}{}
	return b
}

func (b *LoginJWTMiddlewareBuilder) CheckLogin() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
