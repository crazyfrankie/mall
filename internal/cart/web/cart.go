package web

import (
	"github.com/gin-gonic/gin"
	"mall/internal/cart/domain"
	"mall/internal/cart/service"
)

type CartHandler struct {
	svc *service.CartService
}

func NewCartHandler(svc *service.CartService) *CartHandler {
	return &CartHandler{
		svc: svc,
	}
}

func (ctl *CartHandler) RegisterRoute(r *gin.Engine) {
	cartGroup := r.Group("cart")
	{
		cartGroup.POST("add")
	}
}

func (ctl *CartHandler) AddItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		type Req struct {
			ProductID uint64
			Quantity  int
			UserID    uint64
		}
		var req Req
		if err := c.Bind(&req); err != nil {
			return
		}

		err := ctl.svc.AddItem(c.Request.Context(), domain.CartItem{
			ProductID: req.ProductID,
			Quantity:  req.Quantity,
			UserID:    req.UserID,
		})

	}
}

func (ctl *CartHandler) AcquireAllItem() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

//
//func (ctl *CartHandler)
