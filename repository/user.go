package repository

import (
	"context"
	"mall/domain"
	"mall/repository/dao"
)

var (
	ErrUserDuplicateName  = dao.ErrUserDuplicateName
	ErrUserDuplicatePhone = dao.ErrUserDuplicatePhone
	ErrRecordNotFound     = dao.ErrRecordNotFound
	ErrUserNotFound       = dao.ErrUserNotFound
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

func (repo *UserRepository) BindInfo(ctx context.Context, user domain.User) error {
	return repo.dao.BindInfo(ctx, user)
}
