package repository

import "mall/internal/product/repository/dao"

type ProductRepository struct {
	dao *dao.ProductDao
}

func NewProductRepository(dao *dao.ProductDao) *ProductRepository {
	return &ProductRepository{
		dao: dao,
	}
}
