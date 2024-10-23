package domain

type CartItem struct {
	ProductID uint64
	Quantity  int
	UserID    uint64
}

type ShoppingCart struct {
	UserID uint64
	Items  []CartItem
}
