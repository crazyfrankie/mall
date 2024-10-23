package repository

import (
	"context"
	"log"
	"mall/internal/product/domain"
	"mall/internal/product/repository/cache"
	"mall/internal/product/repository/dao"
)

var (
	ErrCategoryDuplicateName = dao.ErrCategoryDuplicateName
	ErrCategoriesNotFound    = dao.ErrCategoriesNotFound
	ErrCategoryNotFound      = dao.ErrCategoryNotFound
	ErrProductNotFound       = dao.ErrProductNotFound
	ErrProductNotOnList      = dao.ErrProductNotOnList
)

type ProductRepository struct {
	dao   *dao.ProductDao
	cache *cache.ProductCache
}

func NewProductRepository(dao *dao.ProductDao, cache *cache.ProductCache) *ProductRepository {
	return &ProductRepository{
		dao:   dao,
		cache: cache,
	}
}

func (repo *ProductRepository) InsertCategory(ctx context.Context, category domain.Category) error {
	return repo.dao.InsertCategory(ctx, category)
}

func (repo *ProductRepository) AcquireAllCategory(ctx context.Context) ([]domain.Category, error) {
	return repo.dao.AcquireAllCategory(ctx)
}

func (repo *ProductRepository) InsertProduct(ctx context.Context, pro domain.Product, attributes []domain.ProductAttribute, imageUrls []string) error {
	return repo.dao.InsertProduct(ctx, pro, attributes, imageUrls)
}

func (repo *ProductRepository) SearchProducts(ctx context.Context, name string) ([]domain.ProductApproximate, error) {
	return repo.dao.SearchProducts(ctx, name)
}

func (repo *ProductRepository) FindProductById(ctx context.Context, id uint64) (domain.ProductDetail, error) {
	detail, err := repo.cache.GetProduct(ctx, id)
	if err == nil {
		return detail, nil
	}
	detail, err = repo.dao.FindProductById(ctx, id)
	if err != nil {
		return domain.ProductDetail{}, err
	}

	// 尝试更新缓存
	err = repo.cache.SetProduct(ctx, id, detail)
	if err != nil {
		log.Printf("缓存更新失败 %v", err.Error())
		// 可选择记录错误或其他处理
	}

	return detail, nil
}

func (repo *ProductRepository) DeleteProductById(ctx context.Context, id uint64) error {
	return repo.dao.DeleteProductById(ctx, id)
}

func (repo *ProductRepository) ProductOn(ctx context.Context, id uint64) error {
	return repo.dao.ProductOn(ctx, id)
}

func (repo *ProductRepository) ProductRemove(ctx context.Context, id uint64) error {
	return repo.dao.ProductRemove(ctx, id)
}
