package dao

import (
	"context"
	"gorm.io/gorm"

	"mall/internal/product/domain"
)

type ProductDao struct {
	db *gorm.DB
}

func NewProductDao(db *gorm.DB) *ProductDao {
	return &ProductDao{
		db: db,
	}
}

func (dao *ProductDao) Insert(ctx context.Context, product domain.Product) error {

}

func (dao *ProductDao) FindProductById(ctx context.Context, id uint64) (domain.Product, error) {

}

func (dao *ProductDao) UnpublishProduct(ctx context.Context, productID uint64) error {
	result := dao.db.WithContext(ctx).Model(&Product{}).Where("id = ?", productID).Update("is_active", false)
	return result.Error
}
