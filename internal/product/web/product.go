package web

import (
	"github.com/gin-gonic/gin"
	"mall/internal/product/service"
)

type ProductHandler struct {
	svc *service.ProductService
}

func NewProductHandler(svc *service.ProductService) *ProductHandler {
	return &ProductHandler{
		svc: svc,
	}
}

func (ctl *ProductHandler) RegisterRoute(r *gin.Engine) {

}

func (ctl *ProductHandler) AddProduct() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func (ctl *ProductHandler) ProductDetail() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

//func (ctl *ProductHandler)
