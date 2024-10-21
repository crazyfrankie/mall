package cache

import (
	"context"
	_ "embed"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
)

var (
	ErrSendTooMany    = errors.New("send too frequency")
	ErrVerifyTooMany  = errors.New("too many verifications")
	ErrCodeStillValid = errors.New("code is still valid")
	ErrKeyConflict    = errors.New("key conflict detected")
	ErrCodeExpired    = errors.New("code has expired")
	ErrCodeNotSet     = errors.New("code is not set")
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
	case -2:
		// 验证码尚未过期
		return ErrCodeStillValid // 新增的错误类型，表示验证码仍然有效
	case -3:
		// 有人误操作导致 key 冲突
		return ErrKeyConflict // 新增的错误类型，表示 key 冲突
	}

	return errors.New("system error") // 保留原来的错误处理
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
		// 频繁错误，可能是恶意攻击，需告警
		return false, ErrVerifyTooMany
	case -2:
		// 验证码已过期
		return false, ErrCodeExpired // 新增的错误类型，表示验证码已过期
	case -3:
		// 验证码未设置或有其他问题
		return false, ErrCodeNotSet // 新增的错误类型，表示验证码未设置
	}

	return false, errors.New("system error") // 保留原来的错误处理
}

func (cache *CodeCache) key(biz, phone string) string {
	return fmt.Sprintf("phone_code:%s:%s", biz, phone)
}
