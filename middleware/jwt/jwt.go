package jwt

import (
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type TokenHandler struct {
	SecretKey []byte
}

type Claim struct {
	jwt.StandardClaims
	sessionId string
}

func NewJwtHandler() *TokenHandler {
	return &TokenHandler{
		SecretKey: []byte("KsS2X1CgFT4bi3BRRIxLk5jjiUBj8wxE"),
	}
}

func (h *TokenHandler) GenerateToken(ctx *gin.Context, sessionId string) error {
	claim := Claim{
		sessionId: sessionId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
		},
	}

	tokenClaim := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	token, err := tokenClaim.SignedString(h.SecretKey)
	ctx.Header("x-jwt-token", token)
	return err
}

func (h *TokenHandler) ExtractToken(ctx *gin.Context) string {
	tokenHeader := ctx.GetHeader("Authorization")
	fmt.Println("Token Header:", tokenHeader) // 打印请求头
	if tokenHeader == "" {
		return ""
	}

	// 分割并检查
	parts := strings.Split(tokenHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "" // 或者返回一个错误
	}

	return parts[1]
}
