package repository

import (
	"context"
	"mall/internal/cart/domain"
	"mall/internal/cart/repository/dao"
)

type CartRepository struct {
	dao *dao.CartDao
}

func NewCartRepository(dao *dao.CartDao) *CartRepository {
	return &CartRepository{
		dao: dao,
	}
}

func (repo *CartRepository) AddToCart(ctx context.Context, item domain.CartItem) error {
	return repo.dao.InsertItem(ctx, repo.domainToDao(item))
}

func (repo *CartRepository) GetCart(ctx context.Context, uid uint64) ([]domain.CartItem, error) {
	var items []domain.CartItem

	carts, err := repo.dao.GetCart(ctx, uid)
	if err != nil {
		return items, err
	}
	for _, item := range carts {
		items = append(items, repo.daoToDomain(item))
	}

	return items, nil
}

func (repo *CartRepository) EmptyCart(ctx context.Context, uid uint64) error {
	return repo.dao.EmptyCart(ctx, uid)
}

func (repo *CartRepository) domainToDao(item domain.CartItem) dao.Cart {
	return dao.Cart{
		UserID:    item.UserID,
		ProductID: item.ProductID,
		Quantity:  item.Quantity,
	}
}

func (repo *CartRepository) daoToDomain(item dao.Cart) domain.CartItem {
	return domain.CartItem{
		ProductID: item.ProductID,
		UserID:    item.UserID,
		Quantity:  item.Quantity,
	}
}
