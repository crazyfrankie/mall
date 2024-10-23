package dao

// CartItem 购物车项模型
type CartItem struct {
	ProductID uint64 `json:"product_id"`
	Quantity  int    `json:"quantity"`
	UserID    uint64 `json:"user_id"` // 或会话ID
}

// ShoppingCart 购物车模型
type ShoppingCart struct {
	UserID uint64     `json:"user_id"`
	Items  []CartItem `json:"items"`
}
