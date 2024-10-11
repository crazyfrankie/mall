package ioc

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"time"
)

func InitDB() *gorm.DB {
	dsn := "root:123456@tcp(localhost:3306)/mall?charset=utf8mb4&parseTime=true&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 表名不加s
		},
		// 可设置外键约束
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		panic("failed to connect database")
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic("failed to connect database")
	}
	// 设置连接池参数
	sqlDB.SetMaxIdleConns(10)                                 // 最大空闲连接数
	sqlDB.SetMaxOpenConns(20)                                 // 最大打开连接数
	sqlDB.SetConnMaxLifetime(time.Duration(60) * time.Minute) // 连接的最大生命周期
	Migrate(db)

	return db
}

func Migrate(db *gorm.DB) {

}
