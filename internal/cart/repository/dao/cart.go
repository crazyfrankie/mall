package dao

import (
	"context"
	"errors"

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

func (dao *CartDao) InsertItem(ctx context.Context, cart Cart) error {
	err := dao.db.WithContext(ctx).
		Model(&Cart{}).
		Where(&Cart{UserID: cart.UserID, ProductID: cart.ProductID}).
		First(&cart).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if cart.ID > 0 {
		return dao.db.WithContext(ctx).
			Model(&Cart{}).
			Where(&Cart{UserID: cart.UserID, ProductID: cart.ProductID}).
			UpdateColumn("quantity", gorm.Expr("quantity+?", cart.Quantity)).Error
	}

	return dao.db.WithContext(ctx).Create(&cart).Error
}

func (dao *CartDao) GetCart(ctx context.Context, uid uint64) ([]Cart, error) {
	var carts []Cart
	err := dao.db.WithContext(ctx).Model(&Cart{}).Where(&Cart{UserID: uid}).Find(&carts).Error
	return carts, err
}

func (dao *CartDao) EmptyCart(ctx context.Context, uid uint64) error {
	if uid == 0 {
		return errors.New("user id is required")
	}

	return dao.db.WithContext(ctx).Delete(&Cart{}, "user_id = ?", uid).Error
}
