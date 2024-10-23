package web

import (
	"errors"
	"github.com/gin-gonic/gin"
	"mall/internal/product/domain"
	"mall/internal/product/service"
	"net/http"
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
	categoryGroup := r.Group("api/categories")
	{
		categoryGroup.POST("/", ctl.AddCategory())
		categoryGroup.GET("/", ctl.GetCategories())
	}

	productGroup := r.Group("api/products")
	{
		productGroup.POST("/", ctl.AddProduct())                      // 添加商品
		productGroup.DELETE("/:id", ctl.AddProduct())                 // 删除商品
		productGroup.POST("/:id/onlist", ctl.ProductOnList())         // 上架商品
		productGroup.POST("/:id/removelist", ctl.ProductRemoveList()) // 下架商品
		productGroup.GET("/search", ctl.SearchProducts())             // 搜索商品
		productGroup.GET("/:id", ctl.GetProductDetail())              // 获取商品详情
	}
}

func (ctl *ProductHandler) AddCategory() gin.HandlerFunc {
	return func(c *gin.Context) {
		type Req struct {
			Name string `json:"name"`
		}
		var req Req
		if err := c.Bind(&req); err != nil {
			return
		}

		err := ctl.svc.AddCategory(c.Request.Context(), domain.Category{
			Name: req.Name,
		})
		switch {
		case errors.Is(err, service.ErrCategoryDuplicateName):
			c.JSON(http.StatusConflict, GetResponse(WithStatus(http.StatusConflict), WithMsg("duplicate category name")))
			return
		case err != nil:
			c.JSON(http.StatusInternalServerError, GetResponse(WithStatus(http.StatusConflict), WithMsg("system error")))
			return
		}

		c.JSON(http.StatusOK, GetResponse(WithStatus(http.StatusOK), WithMsg("add category successfully")))
	}
}

func (ctl *ProductHandler) GetCategories() gin.HandlerFunc {
	return func(c *gin.Context) {
		categories, err := ctl.svc.GetCategories(c.Request.Context())
		switch {
		case errors.Is(err, service.ErrCategoriesNotFound):
			c.JSON(http.StatusNotFound, GetResponse(WithStatus(http.StatusNotFound), WithMsg("no categories be found")))
			return
		case err != nil:
			c.JSON(http.StatusInternalServerError, GetResponse(WithStatus(http.StatusInternalServerError), WithMsg("system error")))
			return
		}

		c.JSON(http.StatusOK, GetResponse(WithStatus(http.StatusOK), WithData(categories)))
	}
}

func (ctl *ProductHandler) AddProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		type Attribute struct {
			Name  string `json:"name"`
			Value string `json:"value"`
		}

		type Req struct {
			Name        string      `json:"name"`
			Description string      `json:"description"`
			Price       float64     `json:"price"`
			Stock       int         `json:"stock"`
			CategoryId  uint64      `json:"categoryId"`
			Images      []string    `json:"images"`
			Attributes  []Attribute `json:"attributes"`
		}
		var req Req
		if err := c.Bind(&req); err != nil {
			return
		}

		var attributes []domain.ProductAttribute
		for _, attr := range req.Attributes {
			attributes = append(attributes, domain.ProductAttribute{
				Name:  attr.Name,
				Value: attr.Value,
			})
		}

		err := ctl.svc.AddProduct(c.Request.Context(), domain.Product{
			Name:        req.Name,
			Description: req.Description,
			Price:       req.Price,
			Stock:       req.Stock,
			CategoryId:  req.CategoryId,
		}, attributes, req.Images)
		switch {
		case errors.Is(err, service.ErrCategoryNotFound):
			c.JSON(http.StatusNotFound, GetResponse(WithStatus(http.StatusNotFound), WithMsg("category not found")))
			return
		case err != nil:
			c.JSON(http.StatusInternalServerError, GetResponse(WithStatus(http.StatusInternalServerError), WithMsg("system error")))
			return
		}

		c.JSON(http.StatusOK, GetResponse(WithStatus(http.StatusOK), WithMsg("add product successfully")))
	}
}

func (ctl *ProductHandler) GetProductDetail() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		product, err := ctl.svc.GetProductDetail(c.Request.Context(), id)
		switch {
		case errors.Is(err, service.ErrProductNotFound):
			c.JSON(http.StatusNotFound, GetResponse(WithStatus(http.StatusNotFound), WithMsg("product not found")))
			return
		case errors.Is(err, service.ErrCategoryNotFound):
			c.JSON(http.StatusNotFound, GetResponse(WithStatus(http.StatusNotFound), WithMsg("category not found")))
			return
		case errors.Is(err, service.ErrProductNotOnList):
			c.JSON(http.StatusNotFound, GetResponse(WithStatus(http.StatusNotFound), WithMsg("product is not on list")))
			return
		case err != nil:
			c.JSON(http.StatusInternalServerError, GetResponse(WithStatus(http.StatusInternalServerError), WithMsg("system error")))
			return
		}

		c.JSON(http.StatusOK, GetResponse(WithStatus(http.StatusOK), WithData(product)))
	}
}

func (ctl *ProductHandler) SearchProducts() gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Query("name")

		products, err := ctl.svc.SearchProducts(c.Request.Context(), name)
		if err != nil {
			c.JSON(http.StatusInternalServerError, GetResponse(WithStatus(http.StatusInternalServerError), WithMsg("system error")))
			return
		}

		c.JSON(http.StatusOK, GetResponse(WithStatus(http.StatusOK), WithData(products)))
	}
}

func (ctl *ProductHandler) DeleteProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		err := ctl.svc.DeleteProduct(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, GetResponse(WithStatus(http.StatusInternalServerError), WithMsg("system error")))
			return
		}

		c.JSON(http.StatusOK, GetResponse(WithStatus(http.StatusOK), WithMsg("delete product successfully")))
	}
}

func (ctl *ProductHandler) ProductOnList() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		err := ctl.svc.ProductOnList(c.Request.Context(), id)
		switch {
		case errors.Is(err, service.ErrProductNotFound):
			c.JSON(http.StatusNotFound, GetResponse(WithStatus(http.StatusNotFound), WithMsg("product not found")))
			return
		case err != nil:
			c.JSON(http.StatusInternalServerError, GetResponse(WithStatus(http.StatusInternalServerError), WithMsg("system error")))
			return
		}

		c.JSON(http.StatusOK, GetResponse(WithStatus(http.StatusOK), WithMsg("product is on list")))
	}
}

func (ctl *ProductHandler) ProductRemoveList() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		err := ctl.svc.ProductOnList(c.Request.Context(), id)
		switch {
		case errors.Is(err, service.ErrProductNotFound):
			c.JSON(http.StatusNotFound, GetResponse(WithStatus(http.StatusNotFound), WithMsg("product not found")))
			return
		case err != nil:
			c.JSON(http.StatusInternalServerError, GetResponse(WithStatus(http.StatusInternalServerError), WithMsg("system error")))
			return
		}

		c.JSON(http.StatusOK, GetResponse(WithStatus(http.StatusOK), WithMsg("product is remove list")))
	}
}
