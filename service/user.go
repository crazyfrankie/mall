package service

import (
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"mall/domain"
	"mall/middleware/jwt"
	"mall/repository"
)

var (
	ErrUserDuplicateName     = repository.ErrUserDuplicateName
	ErrUserDuplicatePhone    = repository.ErrUserDuplicatePhone
	ErrRecordNotFound        = repository.ErrRecordNotFound
	ErrUserNotFound          = repository.ErrUserNotFound
	ErrInvalidUserOrPassword = errors.New("username or password error")
)

type UserService struct {
	repo    *repository.UserRepository
	sessHdl *jwt.RedisSession
}

func NewUserService(repo *repository.UserRepository, sessHdl *jwt.RedisSession) *UserService {
	return &UserService{
		repo:    repo,
		sessHdl: sessHdl,
	}
}

func (svc *UserService) CheckPhone(ctx context.Context, phone string) error {
	return svc.repo.CheckPhone(ctx, phone)
}

func (svc *UserService) CreateUser(ctx context.Context, phone string) error {
	return svc.repo.Insert(ctx, phone)
}

func (svc *UserService) SetSession(ctx context.Context, phone string) (string, error) {
	// 查找用户
	user, err := svc.repo.FindByPhone(ctx, phone)
	if err != nil {
		return "", err
	}
	var ssid string
	ssid, err = svc.sessHdl.CreateSession(ctx, user)
	if err != nil {
		return "", err
	}

	return ssid, nil
}

func (svc *UserService) BindInfo(ctx context.Context, user domain.User) error {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashPassword)

	return svc.repo.BindInfo(ctx, user)
}

func (svc *UserService) NameLogin(ctx context.Context, user domain.User) (string, error) {
	u, err := svc.repo.FindByName(ctx, user.Name)
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(user.Password))
	if err != nil {
		return "", err
	}

	return u.Phone, nil
}
