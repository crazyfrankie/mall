package domain

type Product struct {
	Id          uint64
	Name        string
	Description string
	Price       float64
	Stock       int
	CreateAt    int64
	UpdateAt    int64
}
