package service

import (
	"context"

	"golang.org/x/crypto/bcrypt"

	"mall/domain"
	"mall/repository"
)

var (
	ErrUserDuplicateName  = repository.ErrUserDuplicateName
	ErrUserDuplicatePhone = repository.ErrUserDuplicatePhone
)

type UserService interface {
	SignUp(ctx context.Context, user domain.User) error
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

func (svc *userService) SignUp(ctx context.Context, user domain.User) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hash)
	return svc.repo.Insert(ctx, user)
}
