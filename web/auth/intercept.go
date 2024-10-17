package auth

import (
	"errors"
	"github.com/gin-gonic/gin"
	jwt2 "mall/middleware/jwt"
	"net/http"
)

type LoginJWTMiddlewareBuilder struct {
	paths      map[string]struct{}
	jwtHdl     *jwt2.TokenHandler
	sessionHdl *jwt2.RedisSession
}

func NewLoginJWTMiddlewareBuilder(jwtHdl *jwt2.TokenHandler, sessionHdl *jwt2.RedisSession) *LoginJWTMiddlewareBuilder {
	return &LoginJWTMiddlewareBuilder{
		paths:      make(map[string]struct{}),
		jwtHdl:     jwtHdl,
		sessionHdl: sessionHdl,
	}
}

func (b *LoginJWTMiddlewareBuilder) IgnorePath(path string) *LoginJWTMiddlewareBuilder {
	b.paths[path] = struct{}{}
	return b
}

func (b *LoginJWTMiddlewareBuilder) CheckLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 如果是 OPTIONS 请求，直接放行
		if c.Request.Method == "OPTIONS" {
			c.Next()
			return
		}

		// 路径校验
		if _, ok := b.paths[c.Request.URL.Path]; ok {
			c.Next()
			return
		}

		// 提取并检查 token
		tokenHeader := b.jwtHdl.ExtractToken(c)

		claims, err := b.jwtHdl.ParseToken(tokenHeader)
		if err != nil {
			code, msg := b.jwtHdl.HandleTokenError(err)
			c.JSON(code, msg)
			c.Abort()
			return
		}

		err = b.sessionHdl.AcquireSession(c.Request.Context(), claims.SessionId)
		if err != nil {
			if errors.Is(err, jwt2.ErrKeyNotFound) {
				c.AbortWithError(http.StatusUnauthorized, errors.New("you need login"))
				return
			}

			c.AbortWithError(http.StatusBadRequest, err)
		}

		// 将解析出来的 Claims 存入上下文
		c.Set("claims", claims)
		// 继续后续的处理
		c.Next()
	}
}
