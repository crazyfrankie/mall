package web

import (
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
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

func (ctl *CartHandler) RegisterRoute(r *server.Hertz) {
	cartGroup := r.Group("cart")
	{
		cartGroup.POST("add")
	}
}

func (ctl *CartHandler) AddItem() app.HandlerFunc {
	return func(c *app.RequestContext) {
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

func (ctl *CartHandler) AcquireAllItem() app.HandlerFunc {
	return func(c *app.RequestContext) {

	}
}

//
//func (ctl *CartHandler)
