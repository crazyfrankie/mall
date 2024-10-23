package dao

import (
	"gorm.io/gorm"
)

type CartDao struct {
	db *gorm.DB
}

func NewCartDao(db *gorm.DB) *CartDao {
	return &CartDao{
		db: db,
	}
}

//func (dao *CartDao) Insert(ctx context.Context, item domain.CartItem) error {
//
//}
