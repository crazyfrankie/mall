package repository

import (
	"context"

	"mall/internal/user/domain"
	"mall/internal/user/repository/dao"
)

var (
	ErrUserDuplicateName  = dao.ErrUserDuplicateName
	ErrUserDuplicatePhone = dao.ErrUserDuplicatePhone
	ErrRecordNotFound     = dao.ErrRecordNotFound
)

type UserRepository struct {
	dao *dao.UserDao
}

func NewUserRepository(dao *dao.UserDao) *UserRepository {
	return &UserRepository{
		dao: dao,
	}
}

func (repo *UserRepository) Insert(ctx context.Context, phone string) error {
	return repo.dao.Insert(ctx, phone)
}

func (repo *UserRepository) CheckPhone(ctx context.Context, phone string) error {
	_, err := repo.dao.FindByPhone(ctx, phone)
	return err
}

func (repo *UserRepository) FindByName(ctx context.Context, name string) (domain.User, error) {
	return repo.dao.FindByName(ctx, name)
}

func (repo *UserRepository) FindByPhone(ctx context.Context, phone string) (domain.User, error) {
	return repo.dao.FindByPhone(ctx, phone)
}

func (repo *UserRepository) UpdatePassword(ctx context.Context, user domain.User) error {
	return repo.dao.UpdatePassword(ctx, user)
}

func (repo *UserRepository) UpdateBirthday(ctx context.Context, user domain.User) error {
	return repo.dao.UpdateBirthday(ctx, user)
}

func (repo *UserRepository) UpdateName(ctx context.Context, user domain.User) error {
	return repo.dao.UpdateName(ctx, user)
}

func (repo *UserRepository) AddAddress(ctx context.Context, addr domain.Address) error {
	return repo.dao.AddAddress(ctx, addr)
}

func (repo *UserRepository) FindAllAddrById(ctx context.Context, userId uint64) ([]domain.Address, error) {
	return repo.dao.FindAllAddrById(ctx, userId)
}

func (repo *UserRepository) DeleteAddress(ctx context.Context, userId uint64) error {
	return repo.dao.DeleteAddrById(ctx, userId)
}

func (repo *UserRepository) UpdateAddress(ctx context.Context, addr domain.Address) (domain.Address, error) {
	return repo.dao.UpdateAddrById(ctx, addr)
}
