package dao

import (
	"context"
	"crypto/rand"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"

	"mall/internal/user/domain"
)

var (
	ErrUserDuplicateName  = errors.New("duplicate name")
	ErrUserDuplicatePhone = errors.New("duplicate phone")
	ErrRecordNotFound     = errors.New("record not found")
)

const charset = "abcdefghijklmnopqrstuvwxyz0123456789"

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

	code, err := dao.GenerateCode()
	if err != nil {
		return err
	}
	user := User{
		Phone: phone,
		Name:  code[:15] + "-" + code[15:],
	}
	now := time.Now().UnixMilli()
	user.CreateAt = now
	user.UpdateAt = now
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

	return dao.UserDaoToDomain(user), nil
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

	return dao.UserDaoToDomain(user), nil
}

func (dao *UserDao) UpdatePassword(ctx context.Context, user domain.User) error {
	// 直接更新用户信息
	result := dao.db.WithContext(ctx).Model(&User{}).Where("id = ?", user.Id).Updates(User{
		Password: user.Password,
	})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("user not found or no updates made")
	}

	return nil
}

func (dao *UserDao) UpdateName(ctx context.Context, user domain.User) error {
	// 直接更新用户信息
	result := dao.db.WithContext(ctx).Model(&User{}).Where("id = ?", user.Id).Updates(User{
		Name: user.Name,
	})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("user not found or no updates made")
	}

	return nil
}

func (dao *UserDao) UpdateBirthday(ctx context.Context, user domain.User) error {
	// 直接更新用户信息
	result := dao.db.WithContext(ctx).Model(&User{}).Where("id = ?", user.Id).Updates(User{
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

func (dao *UserDao) AddAddress(ctx context.Context, addr domain.Address) error {
	address := Address{
		UserId:    addr.UserId,
		Street:    addr.Street,
		City:      addr.City,
		State:     addr.State,
		ZipCode:   addr.ZipCode,
		Country:   addr.Country,
		IsDefault: addr.IsDefault,
	}
	result := dao.db.WithContext(ctx).Create(&address)
	return result.Error
}

func (dao *UserDao) FindAllAddrById(ctx context.Context, userId uint64) ([]domain.Address, error) {
	var modelAddresses []Address
	result := dao.db.WithContext(ctx).Where("user_id = ?", userId).Find(&modelAddresses)

	if result.Error != nil {
		return nil, result.Error
	}

	var addresses []domain.Address
	for _, addr := range modelAddresses {
		addresses = append(addresses, dao.AddrDaoToDomain(addr))
	}
	return addresses, nil
}

func (dao *UserDao) DeleteAddrById(ctx context.Context, id uint64) error {
	result := dao.db.WithContext(ctx).Where("user_id ? =", id).Delete(&Address{})
	return result.Error
}

func (dao *UserDao) UpdateAddrById(ctx context.Context, addr domain.Address) (domain.Address, error) {
	updates := Address{}
	if addr.Street != "" {
		updates.Street = addr.Street
	}
	if addr.City != "" {
		updates.City = addr.City
	}
	if addr.State != "" {
		updates.State = addr.State
	}
	if addr.ZipCode != "" {
		updates.ZipCode = addr.ZipCode
	}
	if addr.Country != "" {
		updates.Country = addr.Country
	}
	updates.IsDefault = addr.IsDefault // 默认可以直接赋值

	result := dao.db.WithContext(ctx).Model(&Address{}).Where("id = ?", addr.Id).Updates(updates)
	if result.Error != nil {
		return domain.Address{}, result.Error
	}

	return addr, nil
}

func (dao *UserDao) UserDomainToDao(user domain.User) User {
	return User{
		Name:     user.Name,
		Password: user.Password,
	}
}

func (dao *UserDao) UserDaoToDomain(user User) domain.User {
	return domain.User{
		Id:       user.Id,
		Name:     user.Name,
		Password: user.Password,
		Birthday: user.Birthday.Time,
		Phone:    user.Phone,
	}
}

func (dao *UserDao) AddrDomainToDao(addr domain.Address) Address {
	return Address{
		UserId:    addr.UserId,
		Street:    addr.Street,
		State:     addr.State,
		City:      addr.City,
		ZipCode:   addr.ZipCode,
		Country:   addr.Country,
		IsDefault: addr.IsDefault,
	}
}

func (dao *UserDao) AddrDaoToDomain(addr Address) domain.Address {
	return domain.Address{
		Id:        addr.Id,
		UserId:    addr.UserId,
		Street:    addr.Street,
		State:     addr.State,
		ZipCode:   addr.ZipCode,
		Country:   addr.Country,
		City:      addr.City,
		IsDefault: addr.IsDefault,
	}
}

func (dao *UserDao) GenerateCode() (string, error) {
	bytes := make([]byte, 20)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	var sb strings.Builder
	sb.Grow(20)

	for _, b := range bytes {
		sb.WriteByte(charset[int(b)%len(charset)])
	}

	return sb.String(), nil
}
