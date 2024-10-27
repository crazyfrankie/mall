package service

import (
	"context"
	"errors"
	"strconv"

	"golang.org/x/crypto/bcrypt"

	"mall/internal/auth/jwt"
	"mall/internal/user/domain"
	"mall/internal/user/repository"
)

var (
	ErrUserDuplicateName     = repository.ErrUserDuplicateName
	ErrRecordNotFound        = repository.ErrRecordNotFound
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

func (svc *UserService) FindOrCreateUser(ctx context.Context, phone string) (domain.User, error) {
	user, err := svc.repo.FindByPhone(ctx, phone)
	if err == nil {
		return user, nil
	}

	if !errors.Is(err, ErrRecordNotFound) {
		return domain.User{}, err
	}

	user, err = svc.repo.Insert(ctx, phone)
	if err != nil {
		return domain.User{}, err
	}

	return user, nil
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

func (svc *UserService) DeleteSession(ctx context.Context, isMerchant bool, id string) error {
	return svc.sessHdl.DeleteSession(ctx, isMerchant, id)
}

func (svc *UserService) ExtendSessionExpiration(ctx context.Context, isMerchant bool, id string) error {
	return svc.sessHdl.ExtendSession(ctx, isMerchant, id)
}

func (svc *UserService) UpdatePassword(ctx context.Context, user domain.User) error {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashPassword)

	return svc.repo.UpdatePassword(ctx, user)
}

func (svc *UserService) UpdateName(ctx context.Context, user domain.User) error {
	return svc.repo.UpdateName(ctx, user)
}

func (svc *UserService) UpdateBirthday(ctx context.Context, user domain.User) error {
	return svc.repo.UpdateBirthday(ctx, user)
}

func (svc *UserService) NameLogin(ctx context.Context, user domain.User) (domain.User, error) {
	u, err := svc.repo.FindByName(ctx, user.Name)
	if err != nil {
		return domain.User{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(user.Password))
	if err != nil {
		return domain.User{}, err
	}

	return u, nil
}

func (svc *UserService) BindAddress(ctx context.Context, address domain.Address) error {
	return svc.repo.AddAddress(ctx, address)
}

func (svc *UserService) AcquireAllAddr(ctx context.Context, userId string) ([]domain.Address, error) {
	id, err := strconv.Atoi(userId)
	if err != nil {
		return nil, err
	}

	return svc.repo.FindAllAddrById(ctx, uint64(id))
}

func (svc *UserService) DeleteAddress(ctx context.Context, userId string) error {
	id, err := strconv.Atoi(userId)
	if err != nil {
		return err
	}

	return svc.repo.DeleteAddress(ctx, uint64(id))
}

func (svc *UserService) UpdateAddress(ctx context.Context, addr domain.Address) (domain.Address, error) {
	return svc.repo.UpdateAddress(ctx, addr)
}
