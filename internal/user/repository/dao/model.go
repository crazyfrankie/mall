package dao

import "database/sql"

type User struct {
	Id         uint64 `gorm:"primaryKey,autoIncrement"`
	Phone      string `gorm:"unique; not null"`
	Name       string `gorm:"unique"`
	Birthday   sql.NullTime
	Password   string
	IsMerchant bool `gorm:"default:false"` // 是否为商家
	CreateAt   int64
	UpdateAt   int64
}

type Address struct {
	Id        uint64 `gorm:"primaryKey,autoIncrement"`
	UserId    uint64 `gorm:"not null;index;foreignKey:UserId;references:Id"`
	Street    string `gorm:"not null"`
	City      string `gorm:"not null"`
	State     string `gorm:"not null"`
	ZipCode   string `gorm:"not null"`
	Country   string `gorm:"not null"`
	IsDefault bool   `gorm:"default:false"`
	Ctime     int64
	Uptime    int64
}
