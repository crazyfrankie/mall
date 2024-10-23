package domain

type Product struct {
	Id          uint64
	Name        string
	Description string
	Price       float64
	Stock       int
	CategoryId  uint64
}

type ProductAttribute struct {
	Id        uint64
	ProductId uint64
	Name      string
	Value     string
}

type Category struct {
	ID       uint64
	Name     string
	ParentID uint64
}

type ProductImage struct {
	Id        uint64
	ProductId uint64
	ImageUrl  string
	IsPrimary bool
}

type ProductDetail struct {
	Product    Product            `json:"product"`
	Category   Category           `json:"category"`
	Images     []ProductImage     `json:"images"`
	Attributes []ProductAttribute `json:"attributes"`
}

type ProductApproximate struct {
	Id          uint64
	Name        string
	Description string
	Price       float64
	ImageUrl    string // 主要图片
}
