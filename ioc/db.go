package ioc

import (
	"time"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	"mall/internal/user/repository/dao"
)

func InitDB() *gorm.DB {
	type Config struct {
		DSN             string `yaml:"dsn"`
		MaxIdleConns    int    `yaml:"maxIdleConns"`
		MaxOpenConns    int    `yaml:"maxOpenConns"`
		ConnMaxLifeTime int    `yaml:"connMaxLifeTime"`
	}

	cfg := Config{
		DSN: "root:123456@tcp(localhost:3306)/mall?charset=utf8mb4&parseTime=true&loc=Local",
	}

	err := viper.UnmarshalKey("mysql", &cfg)
	if err != nil {
		panic(err)
	}
	db, err := gorm.Open(mysql.Open(cfg.DSN), &gorm.Config{
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
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)                                        // 最大空闲连接数
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)                                        // 最大打开连接数
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifeTime) * time.Minute * 3) // 连接的最大生命周期
	Migrate(db)

	return db
}

func Migrate(db *gorm.DB) {
	err := db.AutoMigrate(&dao.User{})
	if err != nil {
		panic(err)
	}
}
