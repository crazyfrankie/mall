package dao

import (
	"context"
	"database/sql"
	"errors"
	"github.com/go-sql-driver/mysql"
	"strings"
	"time"

	"gorm.io/gorm"

	"mall/domain"
)

var (
	ErrUserDuplicateName  = errors.New("duplicate name")
	ErrUserDuplicatePhone = errors.New("duplicate phone")
)

type User struct {
	Id       uint64 `gorm:"primaryKey,autoIncrement"`
	Name     string `gorm:"unique;not null"`
	Password string
	Phone    string `gorm:"unique"`
	Ctime    int64
	Uptime   int64
	Birthday sql.NullTime
}

type UserDao interface {
	Insert(ctx context.Context, user domain.User) error
	//Delete(ctx context.Context, id uint64) error
	//Update(ctx context.Context, user domain.User) error
	//FindById(ctx context.Context, id uint64) (domain.User, error)
	//FindOrCreateByPhone(ctx context.Context, phone string) (domain.User, error)

	DomainToDao(user domain.User) User
	DaoToDomain(user User) domain.User
}

type userDao struct {
	db *gorm.DB
}

func NewUserDao(db *gorm.DB) UserDao {
	return &userDao{
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

func (dao *userDao) Insert(ctx context.Context, u domain.User) error {
	user := dao.DomainToDao(u)
	now := time.Now().UnixMilli()
	user.Ctime = now
	user.Uptime = now
	err := dao.db.WithContext(ctx).Create(&user).Error
	return handleError(err)
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

func (dao *userDao) DomainToDao(user domain.User) User {
	return User{
		Name:     user.Name,
		Password: user.Password,
	}
}

func (dao *userDao) DaoToDomain(user User) domain.User {
	return domain.User{
		Id:       user.Id,
		Name:     user.Name,
		Password: user.Password,
		Birthday: user.Birthday.Time,
	}
}
