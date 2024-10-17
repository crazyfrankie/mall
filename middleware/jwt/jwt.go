package jwt

import (
	"errors"
	"net/http"
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
	SessionId string
}

func NewJwtHandler() *TokenHandler {
	return &TokenHandler{
		SecretKey: []byte("KsS2X1CgFT4bi3BRRIxLk5jjiUBj8wxE"),
	}
}

func (h *TokenHandler) GenerateToken(ctx *gin.Context, sessionId string) error {
	claim := Claim{
		SessionId: sessionId,
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

func (h *TokenHandler) ParseToken(token string) (*Claim, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claim{}, func(token *jwt.Token) (interface{}, error) {
		return h.SecretKey, nil
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
		if claims, ok := tokenClaims.Claims.(*Claim); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, ErrTokenInvalid
}

func (h *TokenHandler) HandleTokenError(err error) (int, string) {
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
