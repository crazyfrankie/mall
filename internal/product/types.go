package product

import "mall/internal/product/web"

type Handler = web.ProductHandler // 暴露出去给 ioc 使用
