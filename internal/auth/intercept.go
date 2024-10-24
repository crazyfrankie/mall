package auth

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"mall/internal/auth/jwt"
)

type TokenEffectiveBuilder struct {
	paths      map[string]struct{}
	jwtHdl     *jwt.TokenHandler
	sessionHdl *jwt.RedisSession
}

func NewTokenEffectiveBuilder(jwtHdl *jwt.TokenHandler, sessionHdl *jwt.RedisSession) *TokenEffectiveBuilder {
	return &TokenEffectiveBuilder{
		paths:      make(map[string]struct{}),
		jwtHdl:     jwtHdl,
		sessionHdl: sessionHdl,
	}
}

func (b *TokenEffectiveBuilder) IgnorePath(path string) *TokenEffectiveBuilder {
	b.paths[path] = struct{}{}
	return b
}

func (b *TokenEffectiveBuilder) Check() gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, ok := b.paths[c.Request.URL.Path]; ok {
			c.Next()
			return
		}

		// 提取并检查 token
		tokenHeader := b.jwtHdl.ExtractToken(c)

		claims, err := b.jwtHdl.ParseToken(tokenHeader)
		if err != nil {
			code, msg := b.jwtHdl.HandleTokenError(err)
			c.AbortWithError(code, errors.New(msg))
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

		// 将解析出来的 Claims 存入上下文
		c.Set("claims", claims)
		c.Next()
	}
}
