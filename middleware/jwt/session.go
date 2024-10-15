package jwt

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"

	"mall/domain"
)

type RedisSession struct {
	cmd redis.Cmdable
}

func NewRedisSession(cmd redis.Cmdable) *RedisSession {
	return &RedisSession{
		cmd: cmd,
	}
}

func (s *RedisSession) CreateSession(ctx context.Context, user domain.User) (string, error) {
	ssid := uuid.New().String()

	sessionData := map[string]interface{}{
		"name":     user.Name,
		"password": user.Password,
	}

	sessionDataBytes, err := json.Marshal(sessionData)
	if err != nil {
		return "", err
	}

	err = s.cmd.Set(ctx, ssid, sessionDataBytes, time.Hour*1).Err()
	if err != nil {
		return "", err
	}

	return ssid, nil
}
