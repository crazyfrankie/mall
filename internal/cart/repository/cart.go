package repository

import "mall/internal/cart/repository/dao"

type CartRepository struct {
	dao *dao.CartDao
}

func NewCartRepository(dao *dao.CartDao) *CartRepository {
	return &CartRepository{
		dao: dao,
	}
}
