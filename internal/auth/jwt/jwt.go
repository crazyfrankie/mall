package jwt

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

var (
	ErrTokenInvalid = errors.New("token is invalid")
	ErrTokenExpired = errors.New("token is expired")
	ErrLoginYet     = errors.New("have not logged in yet")
)

type TokenHandler struct {
	SecretKey []byte
}

type Claim struct {
	jwt.StandardClaims
	SessionId  string
	IsMerchant bool
}

func NewJwtHandler() *TokenHandler {
	return &TokenHandler{
		SecretKey: []byte("KsS2X1CgFT4bi3BRRIxLk5jjiUBj8wxE"),
	}
}

func (h *TokenHandler) GenerateToken(ctx *gin.Context, sessionId string, isMerchant bool) error {
	claim := Claim{
		SessionId:  sessionId,
		IsMerchant: isMerchant,
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
	// 分割并检查
	parts := strings.Split(tokenHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "" // 或者返回一个错误
	}

	return parts[1] // 返回 token
}

func (h *TokenHandler) ParseToken(token string) (*Claim, error) {
	// 直接解析 token，其有效性已在网关中确认
	tokenClaims, err := jwt.ParseWithClaims(token, &Claim{}, func(token *jwt.Token) (interface{}, error) {
		return h.SecretKey, nil
	})

	if err != nil {
		// 处理解析错误
		return nil, err // 这里可以返回解析错误，表示解析失败
	}

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claim); ok {
			return claims, nil
		}
	}

	return nil, ErrTokenInvalid // 处理无效 claims
}

func (h *TokenHandler) HandleTokenError(err error) (int, string) {
	var code int
	var msg string
	switch {
	case errors.Is(err, ErrTokenInvalid):
		code = http.StatusUnauthorized
		msg = "token is invalid"
	default:
		code = http.StatusInternalServerError
		msg = "parse Token failed"
	}
	return code, msg
}
