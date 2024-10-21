package dao

type Product struct {
	Id          uint64 `gorm:"primaryKey,autoIncrement"`
	Name        string `gorm:"not null"`
	Description string
	Price       float64 `gorm:"not null"`
	Stock       int     `gorm:"not null"`
	IsActive    bool    `gorm:"default:true"`
	CreateAt    int64
	UpdateAt    int64
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
