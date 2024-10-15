package jwt

import (
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
	if tokenHeader == "" {
		return ""
	}

	return tokenHeader
}
