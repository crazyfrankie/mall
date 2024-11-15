package service

import (
	"context"
	"fmt"
	"mall/internal/cart/domain"
	"mall/internal/cart/repository"
	repository2 "mall/internal/product/repository"
	"net/http"
)

type CartService struct {
	cartRepo    *repository.CartRepository
	productRepo *repository2.ProductRepository
}

type CustomErr struct {
	Code int
	Msg  string
}

func NewCartService(cartRepo *repository.CartRepository, productRepo *repository2.ProductRepository) *CartService {
	return &CartService{
		cartRepo:    cartRepo,
		productRepo: productRepo,
	}
}

func (svc *CartService) AddCart(ctx context.Context, item domain.CartItem) *CustomErr {
	product, err := svc.productRepo.FindProductById(ctx, item.ProductID)
	if err != nil {
		return &CustomErr{Code: http.StatusInternalServerError, Msg: err.Error()}
	}
	if product.Product.Id == 0 {
		return &CustomErr{Code: http.StatusNotFound, Msg: "商品不存在"}
	}

	err = svc.cartRepo.AddToCart(ctx, item)
	if err != nil {
		return &CustomErr{Code: http.StatusInternalServerError, Msg: fmt.Sprintf("用户:%d:添加商品到购物车失败", item.UserID)}
	}

	return nil
}

func (svc *CartService) EmptyCart(ctx context.Context, uid uint64) *CustomErr {
	err := svc.cartRepo.EmptyCart(ctx, uid)
	if err != nil {
		return &CustomErr{Code: http.StatusInternalServerError, Msg: fmt.Sprintf("用户:%d:清空购物车失败", uid)}
	}

	return nil
}

func (svc *CartService) GetCart(ctx context.Context, uid uint64) ([]domain.CartItem, *CustomErr) {
	items, err := svc.cartRepo.GetCart(ctx, uid)
	if err != nil {
		return items, &CustomErr{Code: http.StatusInternalServerError, Msg: fmt.Sprintf("用户:%d:获取购物车失败", uid)}
	}

	return items, nil
}
