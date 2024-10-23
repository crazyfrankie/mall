package auth

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"mall/internal/auth/jwt"
)

type AuthorityJWTMiddlewareBuilder struct {
	merchantPaths map[string]struct{}
	ignorePaths   map[string]struct{}
	jwtHdl        *jwt.TokenHandler
	sessionHdl    *jwt.RedisSession
}

func NewAuthorityMiddlewareBuilder(jwtHdl *jwt.TokenHandler, sessionHdl *jwt.RedisSession) *AuthorityJWTMiddlewareBuilder {
	return &AuthorityJWTMiddlewareBuilder{
		merchantPaths: make(map[string]struct{}),
		ignorePaths:   make(map[string]struct{}),
		jwtHdl:        jwtHdl,
		sessionHdl:    sessionHdl,
	}
}

func (b *AuthorityJWTMiddlewareBuilder) IgnorePath(path string) *AuthorityJWTMiddlewareBuilder {
	b.ignorePaths[path] = struct{}{}
	return b
}

func (b *AuthorityJWTMiddlewareBuilder) MerchantPath(path string) *AuthorityJWTMiddlewareBuilder {
	b.merchantPaths[path] = struct{}{}
	return b
}

func (b *AuthorityJWTMiddlewareBuilder) CheckLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 路径校验
		if _, ok := b.ignorePaths[c.Request.URL.Path]; ok {
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

		// 检测 ssid 是否有效
		err = b.sessionHdl.AcquireSession(c.Request.Context(), claims.SessionId)
		if err != nil {
			if errors.Is(err, jwt.ErrKeyNotFound) {
				c.AbortWithError(http.StatusUnauthorized, errors.New("you need login"))
				return
			}

			c.AbortWithError(http.StatusBadRequest, err)
		}

		// 刷新 ssid 有效时间
		err = b.sessionHdl.ExtendSession(c.Request.Context(), claims.SessionId)
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		// 检查 JWT 是否快过期，若是，则生成新 JWT
		if time.Until(time.Unix(claims.ExpiresAt, 0)) < time.Minute*5 {
			err := b.jwtHdl.GenerateToken(c, claims.SessionId, claims.IsMerchant)
			if err != nil {
				c.AbortWithStatus(http.StatusBadRequest)
				return
			}
		}

		// 检查是否为商家
		if _, isMerchantPath := b.merchantPaths[c.Request.URL.Path]; isMerchantPath && !claims.IsMerchant {
			c.AbortWithStatusJSON(http.StatusForbidden, "access denied for non-merchants")
			return
		}

		// 将解析出来的 Claims 存入上下文
		c.Set("claims", claims)
		c.Next()
	}
}
