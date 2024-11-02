package dao

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	"mall/internal/product/domain"
)

var (
	ErrCategoryNotFound      = errors.New("category does not exist")
	ErrCategoryDuplicateName = errors.New("duplicate category")
	ErrCategoriesNotFound    = errors.New("empty categories found")
	ErrProductNotFound       = errors.New("product not found")
	ErrProductNotOnList      = errors.New("product is not on list")
)

type ProductDao struct {
	db *gorm.DB
}

func NewProductDao(db *gorm.DB) *ProductDao {
	return &ProductDao{
		db: db,
	}
}

func (dao *ProductDao) InsertCategory(ctx context.Context, category domain.Category) error {
	cg := Category{
		Name: category.Name,
	}
	now := time.Now().UnixMilli()
	cg.CreatedAt = now
	cg.UpdatedAt = now
	if err := dao.db.WithContext(ctx).Create(&cg).Error; err != nil {
		return err
	}
	return nil
}

func (dao *ProductDao) AcquireAllCategory(ctx context.Context) ([]domain.Category, error) {
	var cgs []Category

	err := dao.db.WithContext(ctx).Model(&Category{}).Find(&cgs).Error
	if err != nil {
		if errors.Is(err, gorm.ErrEmptySlice) {
			return []domain.Category{}, ErrCategoriesNotFound
		}

		return []domain.Category{}, err
	}

	var categories []domain.Category
	for _, cg := range cgs {
		categories = append(categories, dao.categoryDaoToDomain(cg))
	}

	return categories, nil
}

func (dao *ProductDao) InsertProductImage(ctx context.Context, productId uint64, imageUrl string, isPrimary bool) error {
	// 插入商品图片
	productImage := ProductImage{
		ProductId: productId,
		ImageUrl:  imageUrl,
		IsPrimary: isPrimary,
	}

	if err := dao.db.WithContext(ctx).Create(&productImage).Error; err != nil {
		return err
	}

	return nil
}

func (dao *ProductDao) InsertProductAttribute(ctx context.Context, productId uint64, attribute domain.ProductAttribute) error {
	// 在这里可以添加额外的逻辑，比如验证属性数据等
	at := ProductAttribute{
		ProductId: productId,
		Name:      attribute.Name,
		Value:     attribute.Value,
	}
	// 插入商品属性
	if err := dao.db.WithContext(ctx).Create(&at).Error; err != nil {
		return err
	}

	return nil
}

func (dao *ProductDao) InsertProduct(ctx context.Context, pro domain.Product, attributes []domain.ProductAttribute, imageUrls []string) error {
	// 检查分类是否存在
	var category Category
	if err := dao.db.WithContext(ctx).Where("id = ?", pro.CategoryId).First(&category).Error; err != nil {
		return ErrCategoryNotFound
	}

	// 插入商品
	if err := dao.db.WithContext(ctx).Create(&pro).Error; err != nil {
		return err
	}

	// 插入关联记录
	productCategory := ProductCategory{
		ProductID:  pro.Id,         // 商品ID
		CategoryID: pro.CategoryId, // 分类ID
	}

	if err := dao.db.WithContext(ctx).Create(&productCategory).Error; err != nil {
		return err
	}

	// 插入商品属性
	for _, attr := range attributes {
		if err := dao.InsertProductAttribute(ctx, pro.Id, attr); err != nil {
			return err
		}
	}

	// 插入商品图片
	for _, imageUrl := range imageUrls {
		isPrimary := imageUrl == imageUrls[0] // 第一张图片设置为主图
		if err := dao.InsertProductImage(ctx, pro.Id, imageUrl, isPrimary); err != nil {
			return err
		}
	}

	return nil
}

func (dao *ProductDao) UpdateProductSock(ctx context.Context, productId uint64, quantity int) error {
	// 查找商品
	var product Product
	if err := dao.db.WithContext(ctx).Where("id = ?", productId).First(&product).Error; err != nil {
		return ErrProductNotFound
	}

	// 更新库存
	product.Stock += quantity
	if err := dao.db.WithContext(ctx).Save(&product).Error; err != nil {
		return err
	}

	return nil
}

func (dao *ProductDao) SearchProducts(ctx context.Context, name string) ([]domain.ProductApproximate, error) {
	var products []Product

	var categoryId uint64
	if name != "" {
		var category Category
		if err := dao.db.WithContext(ctx).Where("name = ?", name).First(&category).Error; err != nil {
			return nil, err // 如果分类未找到，可以返回错误
		}
		categoryId = category.ID
	}

	// 构建查询
	query := dao.db.WithContext(ctx).Model(&Product{})
	if name != "" {
		query = query.Where("name LIKE ? AND is_active = ?", "%"+name+"%", true)
	}
	if categoryId > 0 {
		query = query.Joins("JOIN product_category pc ON pc.product_id = products.id").
			Where("pc.category_id = ? AND products.is_active = ?", categoryId, true)
	}

	if err := query.Find(&products).Error; err != nil {
		return nil, err
	}

	var domainProducts []domain.ProductApproximate
	for _, p := range products {
		domainProducts = append(domainProducts, domain.ProductApproximate{
			Id:          p.Id,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
		})
	}

	return domainProducts, nil
}

func (dao *ProductDao) FindProductById(ctx context.Context, id uint64) (domain.ProductDetail, error) {
	type ProductDetailResult struct {
		Product    Product
		Category   Category
		Attributes []ProductAttribute
		Images     []ProductImage
	}

	var res ProductDetailResult
	// 查找商品，包括是否上架
	if err := dao.db.WithContext(ctx).Where("id = ? AND is_active = ?", id, true).First(&res.Product).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ProductDetail{}, ErrProductNotOnList
		}
		return domain.ProductDetail{}, ErrProductNotFound
	}

	// 查找分类
	if err := dao.db.WithContext(ctx).Where("id = (SELECT category_id FROM product_category WHERE product_id = ?)", id).First(&res.Category).Error; err != nil {
		return domain.ProductDetail{}, ErrCategoryNotFound
	}

	// 查找商品图片
	if err := dao.db.WithContext(ctx).Where("product_id = ?", id).Find(&res.Images).Error; err != nil {
		return domain.ProductDetail{}, err
	}

	// 查找商品属性
	if err := dao.db.WithContext(ctx).Where("product_id = ?", id).Find(&res.Attributes).Error; err != nil {
		return domain.ProductDetail{}, err
	}

	return domain.ProductDetail{
		Product:    dao.productDaoToDomain(res.Product),
		Category:   dao.categoryDaoToDomain(res.Category),
		Images:     dao.imageDaoToDomain(res.Images),
		Attributes: dao.attributeDaoToDomain(res.Attributes),
	}, nil
}

func (dao *ProductDao) DeleteProductById(ctx context.Context, id uint64) error {
	err := dao.db.WithContext(ctx).Where("id = ?", id).Delete(&Product{}).Error
	return err
}

func (dao *ProductDao) ProductOn(ctx context.Context, id uint64) error {
	err := dao.db.WithContext(ctx).Where("id = ?", id).Update("is_active", true).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrProductNotFound
		}
		return err
	}

	return nil
}

func (dao *ProductDao) ProductRemove(ctx context.Context, id uint64) error {
	err := dao.db.WithContext(ctx).Where("id = ?", id).Update("is_active", false).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrProductNotFound
		}
		return err
	}

	return nil
}

//func (dao *ProductDao) FindProductById(ctx context.Context, id uint64) (domain.ProductDetail, error) {
//	type ProductDetailResult struct {
//		Product  Product
//		Category Category
//		Images   []ProductImage
//	}
//
//	var res ProductDetailResult
//	if err := dao.db.WithContext(ctx).
//		Table("products").
//		Select("products.*, (SELECT category_id FROM product_category WHERE product_id = products.id) AS category_id").
//		Joins("LEFT JOIN product_category pc ON pc.product_id = products.id").
//		Joins("LEFT JOIN product_images pi ON pi.product_id = products.id").
//		Where("products.id = ? AND products.is_active = ?", id, true).
//		Group("products.id").
//		Scan(&res).Error; err != nil {
//		return domain.ProductDetail{}, ErrProductNotFound
//	}
//
//	// 查找分类和图片
//	if err := dao.db.WithContext(ctx).Where("id = ?", res.Category.ID).First(&res.Category).Error; err != nil {
//		return domain.ProductDetail{}, ErrCategoryNotFound
//	}
//
//	if err := dao.db.WithContext(ctx).Where("product_id = ?", id).Find(&res.Images).Error; err != nil {
//		return domain.ProductDetail{}, err
//	}
//
//	return domain.ProductDetail{
//		Product:  dao.productDaoToDomain(res.Product),
//		Category: dao.categoryDaoToDomain(res.Category),
//		Images:   dao.imageDaoToDomain(res.Images),
//	}, nil
//}

func (dao *ProductDao) categoryDomainToDao(category domain.Category) Category {
	return Category{
		ID:       category.ID,
		Name:     category.Name,
		ParentID: category.ParentID,
	}
}

func (dao *ProductDao) categoryDaoToDomain(category Category) domain.Category {
	return domain.Category{
		ID:       category.ID,
		Name:     category.Name,
		ParentID: category.ParentID,
	}
}

func (dao *ProductDao) productDaoToDomain(product Product) domain.Product {
	return domain.Product{
		Id:          product.Id,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Stock:       product.Stock,
	}
}

func (dao *ProductDao) imageDaoToDomain(images []ProductImage) []domain.ProductImage {
	var domainImages []domain.ProductImage
	for _, img := range images {
		domainImages = append(domainImages, domain.ProductImage{
			Id:        img.Id,
			ProductId: img.ProductId,
			ImageUrl:  img.ImageUrl,
			IsPrimary: img.IsPrimary,
		})
	}
	return domainImages
}

func (dao *ProductDao) attributeDaoToDomain(attributes []ProductAttribute) []domain.ProductAttribute {
	var ats []domain.ProductAttribute
	for _, at := range attributes {
		ats = append(ats, domain.ProductAttribute{
			Id:        at.Id,
			ProductId: at.ProductId,
			Name:      at.Name,
			Value:     at.Value,
		})
	}

	return ats
}
