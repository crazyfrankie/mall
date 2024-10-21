package repository

import (
	"context"

	"mall/internal/user/repository/cache"
)

var (
	ErrSendTooMany   = cache.ErrSendTooMany
	ErrVerifyTooMany = cache.ErrVerifyTooMany
)

type CodeRepository struct {
	cache *cache.CodeCache
}

func NewCodeRepository(cache *cache.CodeCache) *CodeRepository {
	return &CodeRepository{
		cache: cache,
	}
}

func (repo *CodeRepository) Store(ctx context.Context, biz, phone, code string) error {
	return repo.cache.Store(ctx, biz, phone, code)
}

func (repo *CodeRepository) Verify(ctx context.Context, biz, phone, code string) (bool, error) {
	return repo.cache.Acquire(ctx, biz, phone, code)
}
