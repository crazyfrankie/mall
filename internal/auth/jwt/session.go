package jwt

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"

	"mall/internal/user/domain"
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

	key := s.key(user.IsMerchant, strconv.Itoa(int(user.Id)))

	err := s.cmd.Set(ctx, key, ssid, time.Hour*24*7).Err()
	if err != nil {
		return "", err
	}

	return ssid, nil
}

func (s *RedisSession) DeleteSession(ctx context.Context, isMerchant bool, id string) error {
	key := s.key(isMerchant, id)
	// 尝试删除 session，返回任何可能的错误
	_, err := s.cmd.Del(ctx, key).Result()
	if err != nil {
		return fmt.Errorf("failed to delete session %s: %w", key, err)
	}
	return nil
}

func (s *RedisSession) AcquireSession(ctx context.Context, isMerchant bool, id string) error {
	key := s.key(isMerchant, id)

	_, err := s.cmd.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return ErrKeyNotFound
		}

		return err
	}

	return nil
}

func (s *RedisSession) ExtendSession(ctx context.Context, isMerchant bool, id string) error {
	key := s.key(isMerchant, id)

	_, err := s.cmd.Expire(ctx, key, time.Hour*1).Result()
	return err
}

func (s *RedisSession) key(isMerchant bool, id string) string {
	if isMerchant {
		return fmt.Sprintf("merchant:%s", id)
	}

	return fmt.Sprintf("user:%s", id)
}
