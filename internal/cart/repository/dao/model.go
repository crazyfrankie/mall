package dao

import "gorm.io/gorm"

// Cart 购物车项模型
type Cart struct {
	gorm.Model
	ProductID uint64 `json:"product_id"`
	Quantity  int    `json:"quantity"`
	UserID    uint64 `json:"user_id"`
}

//// ShoppingCart 购物车模型
//type ShoppingCart struct {
//	UserID uint64     `json:"user_id"`
//	Items  []CartItem `json:"items"`
//}
