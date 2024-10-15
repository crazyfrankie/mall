package dao

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"

	"mall/domain"
)

var (
	ErrUserDuplicateName  = errors.New("duplicate name")
	ErrUserDuplicatePhone = errors.New("duplicate phone")
	ErrRecordNotFound     = errors.New("record not found")
	ErrUserNotFound       = errors.New("user not found")
)

type User struct {
	Id       uint64         `gorm:"primaryKey,autoIncrement"`
	Phone    string         `gorm:"unique; not null"`
	Name     sql.NullString `gorm:"unique"`
	Birthday sql.NullTime
	Password string
	Ctime    int64
	Uptime   int64
}

type UserDao struct {
	db *gorm.DB
}

func NewUserDao(db *gorm.DB) *UserDao {
	return &UserDao{
		db: db,
	}
}

func handleError(err error) error {
	var mysqlErr *mysql.MySQLError
	const uniqueConflictErrNo uint16 = 1062

	if errors.As(err, &mysqlErr) && mysqlErr.Number == uniqueConflictErrNo {
		if strings.Contains(mysqlErr.Message, "name") {
			return ErrUserDuplicateName
		} else if strings.Contains(mysqlErr.Message, "phone") {
			return ErrUserDuplicatePhone
		}
	}
	return err
}

func (dao *UserDao) Insert(ctx context.Context, phone string) error {
	// 开启事务
	tx := dao.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	var existingUser User
	if err := tx.WithContext(ctx).Where("phone = ?", phone).First(&existingUser).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			tx.Rollback() // 查询出错则回滚事务
			return err
		}
	} else {
		tx.Rollback() // 如果用户已存在则回滚事务
		return ErrUserDuplicatePhone
	}

	user := User{
		Phone: phone,
	}
	now := time.Now().UnixMilli()
	user.Ctime = now
	user.Uptime = now
	if err := tx.WithContext(ctx).Create(&user).Error; err != nil {
		tx.Rollback() // 插入失败则回滚事务
		return handleError(err)
	}
	return tx.Commit().Error
}

func (dao *UserDao) FindByName(ctx context.Context, name string) (domain.User, error) {
	var user User

	result := dao.db.WithContext(ctx).Where("name = ?", name).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return domain.User{}, ErrRecordNotFound
		}
		return domain.User{}, result.Error
	}

	return dao.DaoToDomain(user), nil
}

func (dao *UserDao) FindByPhone(ctx context.Context, phone string) (domain.User, error) {
	var user User
	result := dao.db.WithContext(ctx).Where("phone = ?", phone).First(&user)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return domain.User{}, ErrRecordNotFound
		}
		return domain.User{}, result.Error
	}

	return dao.DaoToDomain(user), nil
}

func (dao *UserDao) BindInfo(ctx context.Context, user domain.User) error {
	// 直接更新用户信息
	result := dao.db.WithContext(ctx).Model(&User{}).Where("phone = ?", user.Phone).Updates(User{
		Name: sql.NullString{
			String: user.Name,
			Valid:  user.Name != "",
		},
		Password: user.Password,
		Birthday: sql.NullTime{
			Time:  user.Birthday,
			Valid: true,
		},
	})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("user not found or no updates made")
	}

	return nil
}

//func (dao *userDao) FindById(ctx context.Context, id uint64) (domain.User, error) {
//
//}
//
//func (dao *userDao) FindOrCreateByPhone(ctx context.Context, phone string) (domain.User, error) {
//
//}
//
//func (dao *userDao) Delete(ctx context.Context, id uint64) error {
//
//}
//
//func (dao *userDao) Update(ctx context.Context, user domain.User) error {
//
//}

func (dao *UserDao) DomainToDao(user domain.User) User {
	return User{
		Name: sql.NullString{
			String: user.Name,
			Valid:  user.Name != "",
		},
		Password: user.Password,
	}
}

func (dao *UserDao) DaoToDomain(user User) domain.User {
	return domain.User{
		Id:       user.Id,
		Name:     user.Name.String,
		Password: user.Password,
		Birthday: user.Birthday.Time,
		Phone:    user.Phone,
	}
}
