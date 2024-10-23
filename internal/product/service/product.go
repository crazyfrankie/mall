package service

import (
	"context"
	"mall/internal/product/domain"
	"mall/internal/product/repository"
	"strconv"
)

var (
	ErrCategoryDuplicateName = repository.ErrCategoryDuplicateName
	ErrCategoriesNotFound    = repository.ErrCategoriesNotFound
	ErrCategoryNotFound      = repository.ErrCategoryNotFound
	ErrProductNotFound       = repository.ErrProductNotFound
	ErrProductNotOnList      = repository.ErrProductNotOnList
)

type ProductService struct {
	repo *repository.ProductRepository
}

func NewProductService(repo *repository.ProductRepository) *ProductService {
	return &ProductService{
		repo: repo,
	}
}

func (svc *ProductService) AddCategory(ctx context.Context, category domain.Category) error {
	return svc.repo.InsertCategory(ctx, category)
}

func (svc *ProductService) GetCategories(ctx context.Context) ([]domain.Category, error) {
	return svc.repo.AcquireAllCategory(ctx)
}

func (svc *ProductService) AddProduct(ctx context.Context, product domain.Product, attributes []domain.ProductAttribute, images []string) error {
	return svc.repo.InsertProduct(ctx, product, attributes, images)
}

func (svc *ProductService) GetProductDetail(ctx context.Context, productId string) (domain.ProductDetail, error) {
	id, err := strconv.Atoi(productId)
	if err != nil {
		return domain.ProductDetail{}, err
	}

	return svc.repo.FindProductById(ctx, uint64(id))
}

func (svc *ProductService) SearchProducts(ctx context.Context, name string) ([]domain.ProductApproximate, error) {
	return svc.repo.SearchProducts(ctx, name)
}

func (svc *ProductService) DeleteProduct(ctx context.Context, id string) error {
	productId, err := strconv.Atoi(id)
	if err != nil {
		return err
	}

	return svc.repo.DeleteProductById(ctx, uint64(productId))
}

func (svc *ProductService) ProductOnList(ctx context.Context, productId string) error {
	id, err := strconv.Atoi(productId)
	if err != nil {
		return err
	}

	return svc.repo.ProductOn(ctx, uint64(id))
}

func (svc *ProductService) ProductRemoveList(ctx context.Context, productId string) error {
	id, err := strconv.Atoi(productId)
	if err != nil {
		return err
	}

	return svc.repo.ProductRemove(ctx, uint64(id))
}
