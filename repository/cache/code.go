package cache

import (
	"context"
	_ "embed"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
)

var (
	ErrSendTooMany   = errors.New("send too frequency")
	ErrVerifyTooMany = errors.New("too many verifications")
)

//go:embed lua/set_code.lua
var luaSetCode string

//go:embed lua/verify_code.lua
var luaVerifyCode string

type CodeCache struct {
	cmd redis.Cmdable
}

func NewCodeCache(cmd redis.Cmdable) *CodeCache {
	return &CodeCache{
		cmd: cmd,
	}
}

func (cache *CodeCache) Store(ctx context.Context, biz, phone, code string) error {
	key := cache.key(biz, phone)
	res, err := cache.cmd.Eval(ctx, luaSetCode, []string{key}, code).Int()
	if err != nil {
		return err
	}

	switch res {
	case 0:
		return nil
	case -1:
		// 发送太频繁
		return ErrSendTooMany
	}

	return errors.New("system error")
}

func (cache *CodeCache) Acquire(ctx context.Context, biz, phone, inputCode string) (bool, error) {
	key := cache.key(biz, phone)

	res, err := cache.cmd.Eval(ctx, luaVerifyCode, []string{key}, inputCode).Int()
	if err != nil {
		return false, err
	}

	switch res {
	case 0:
		return true, nil
	case -1:
		// 如果频繁出这个错误代表有人搞你 需要告警
		return false, ErrVerifyTooMany
	}

	return false, errors.New("system error")
}

func (cache *CodeCache) key(biz, phone string) string {
	return fmt.Sprintf("phone_code:%s:%s", biz, phone)
}
