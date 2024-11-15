package web

import (
	"github.com/gin-gonic/gin"
	"mall/internal/cart/domain"
	"mall/internal/cart/service"
	"net/http"
)

type CartHandler struct {
	svc *service.CartService
}

func NewCartHandler(svc *service.CartService) *CartHandler {
	return &CartHandler{
		svc: svc,
	}
}

func (ctl *CartHandler) AddCartItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		type Req struct {
			ProductID uint64 `json:"productI_id"`
			UserID    uint64 `json:"user_id"`
			Quantity  int    `json:"quantity"`
		}
		var req Req
		if err := c.Bind(&req); err != nil {
			return
		}

		err := ctl.svc.AddCart(c.Request.Context(), domain.CartItem{
			UserID:    req.UserID,
			ProductID: req.ProductID,
			Quantity:  req.Quantity,
		})
		if err != nil {
			c.JSON(err.Code, err.Msg)
			return
		}

		c.JSON(http.StatusOK, "添加购物车成功")
	}
}

func (ctl *CartHandler) GetCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		type Req struct {
			UserID uint64 `json:"user_id"`
		}
		var req Req
		if err := c.Bind(&req); err != nil {
			return
		}

		carts, err := ctl.svc.GetCart(c.Request.Context(), req.UserID)
		if err != nil {
			c.JSON(err.Code, err.Msg)
			return
		}

		c.JSON(http.StatusOK, carts)
	}
}

func (ctl *CartHandler) EmptyCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		type Req struct {
			UserID uint64 `json:"user_id"`
		}
		var req Req
		if err := c.Bind(&req); err != nil {
			return
		}

		err := ctl.svc.EmptyCart(c.Request.Context(), req.UserID)
		if err != nil {
			c.JSON(err.Code, err.Msg)
			return
		}

		c.JSON(http.StatusOK, "清空购物车成功")
	}
}
