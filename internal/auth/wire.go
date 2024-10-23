package auth

import (
	"github.com/google/wire"
	"mall/internal/auth/jwt"
)

var JWTSet = wire.NewSet(
	jwt.NewJwtHandler,
	jwt.NewRedisSession,
)
