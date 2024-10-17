package jwt

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"

	"mall/domain"
)

var (
	ErrKeyNotFound = errors.New("key this not found")
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

	err = s.cmd.Set(ctx, ssid, sessionDataBytes, time.Hour*24*7).Err()
	if err != nil {
		return "", err
	}

	return ssid, nil
}

func (s *RedisSession) DeleteSession(ctx context.Context, ssid string) error {
	// 尝试删除 session，返回任何可能的错误
	_, err := s.cmd.Del(ctx, ssid).Result()
	if err != nil {
		return fmt.Errorf("failed to delete session %s: %w", ssid, err)
	}
	return nil
}

func (s *RedisSession) AcquireSession(ctx context.Context, ssid string) error {
	_, err := s.cmd.Get(ctx, ssid).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return ErrKeyNotFound
		}

		return err
	}

	return nil
}
