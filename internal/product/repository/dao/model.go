package dao

type Product struct {
	Id          uint64 `gorm:"primaryKey,autoIncrement"`
	Name        string `gorm:"not null"`
	Description string
	Price       float64            `gorm:"not null"`
	Stock       int                `gorm:"not null"`
	IsActive    bool               `gorm:"default:true,index"`
	Quantity    int                `gorm:"not null"`
	Attributes  []ProductAttribute `gorm:"type:json"` // 存储为 JSON 类型
	CreateAt    int64
	UpdateAt    int64
}

type ProductImage struct {
	Id        uint64 `gorm:"primaryKey,autoIncrement"`
	ProductId uint64 `gorm:"not null"`      // 外键，关联到 Product 表
	ImageUrl  string `gorm:"not null"`      // 图片的 URL
	IsPrimary bool   `gorm:"default:false"` // 是否是主图
}

type ProductAttribute struct {
	Id        uint64 `gorm:"primaryKey,autoIncrement"`
	ProductId uint64 // 外键
	Name      string `json:"name"`
	Value     string `json:"value"`
}

type Category struct {
	ID        uint64 `gorm:"primaryKey,autoIncrement"`
	Name      string `gorm:"unique;not null"`
	ParentID  uint64 `gorm:"default:0"` // 支持父子分类
	CreatedAt int64
	UpdatedAt int64
}

type ProductCategory struct {
	ProductID  uint64 `gorm:"primaryKey,autoIncrement:false;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	CategoryID uint64 `gorm:"primaryKey,autoIncrement:false;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
