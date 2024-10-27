package service

import (
	"mall/internal/cart/repository"
)

type CartService struct {
	repo *repository.CartRepository
}

func NewCartService(repo *repository.CartRepository) *CartService {
	return &CartService{
		repo: repo,
	}
}

//func (svc *CartService) AddItem(ctx context.Context, item domain.CartItem) error {
//
//}
