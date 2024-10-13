package repository

import (
	"context"

	"mall/domain"
	"mall/repository/dao"
)

var (
	ErrUserDuplicateName  = dao.ErrUserDuplicateName
	ErrUserDuplicatePhone = dao.ErrUserDuplicatePhone
)

type UserRepository interface {
	Insert(ctx context.Context, user domain.User) error
}

type userRepository struct {
	dao dao.UserDao
}

func NewUserRepository(dao dao.UserDao) UserRepository {
	return &userRepository{
		dao: dao,
	}
}

func (repo *userRepository) Insert(ctx context.Context, user domain.User) error {
	return repo.dao.Insert(ctx, user)
}
