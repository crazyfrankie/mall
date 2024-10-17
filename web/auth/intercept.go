package auth

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"

	jwt2 "mall/middleware/jwt"
)

var (
	ErrTokenInvalid = errors.New("token is invalid")
	ErrTokenExpired = errors.New("token is expired")
	ErrLoginYet     = errors.New("have not logged in yet")
)

type LoginJWTMiddlewareBuilder struct {
	paths  map[string]struct{}
	jwtHdl *jwt2.TokenHandler
}

func NewLoginJWTMiddlewareBuilder(jwtHdl *jwt2.TokenHandler) *LoginJWTMiddlewareBuilder {
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

		claims, err := ParseToken(tokenHeader, string(b.jwtHdl.SecretKey))
		if err != nil {
			code, msg := handleTokenError(err)
			c.JSON(code, msg)
			c.Abort()
			return
		}

		// 将解析出来的 Claims 存入上下文
		c.Set("claims", claims)
		// 继续后续的处理
		c.Next()
	}
}

func ParseToken(token string, key string) (*jwt2.Claim, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &jwt2.Claim{}, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})
	if err != nil {
		var ve *jwt.ValidationError
		if errors.As(err, &ve) {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, ErrTokenInvalid
			} else if ve.Errors&(jwt.ValidationErrorExpired) != 0 {
				return nil, ErrTokenExpired
			} else if ve.Errors&(jwt.ValidationErrorNotValidYet) != 0 {
				return nil, ErrLoginYet
			}
		}
		return nil, err
	}
	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*jwt2.Claim); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, ErrTokenInvalid
}

func handleTokenError(err error) (int, string) {
	var code int
	var msg string
	switch {
	case errors.Is(err, ErrTokenExpired):
		code = http.StatusUnauthorized
		msg = "token is expired"
	case errors.Is(err, ErrTokenInvalid):
		code = http.StatusUnauthorized
		msg = "token is invalid"
	case errors.Is(err, ErrLoginYet):
		code = http.StatusUnauthorized
		msg = "have not logged in yet"
	default:
		code = http.StatusInternalServerError
		msg = "parse Token failed"
	}
	return code, msg
}
