package service

import (
	"greenenvironment/constant"
	"greenenvironment/features/impacts"
	"greenenvironment/features/products"

	"github.com/google/uuid"
)

type ProductService struct {
	productRepo products.ProductRepositoryInterface
	impactRepo  impacts.ImpactRepositoryInterface
}

func NewProductService(pr products.ProductRepositoryInterface, ir impacts.ImpactRepositoryInterface) products.ProductServiceInterface {
	return &ProductService{productRepo: pr, impactRepo: ir}
}

func (ps *ProductService) Create(product products.Product) error {
	product.ID = uuid.New().String()
	for i, impact := range product.ImpactCategories {
		data, _ := ps.impactRepo.GetByID(impact.ImpactCategoryID)
		if data.ID == "" {
			return constant.ErrCreateProduct
		}
		impact.ID = uuid.New().String()
		product.ImpactCategories[i] = impact
	}

	for i, image := range product.Images {
		image.ID = uuid.New().String()
		product.Images[i] = image
	}

	return ps.productRepo.Create(product)
}

func (ps *ProductService) GetAllByPage(page int, search string, sort string) ([]products.Product, int, error) {
	products, total, err := ps.productRepo.GetAllByPage(page, search, sort)
	if err != nil {
		return nil, 0, err
	}
	if page > total {
		return nil, 0, constant.ErrPageInvalid
	}

	return products, total, nil
}

func (ps *ProductService) GetByCategory(category string, page int, search string, sort string) ([]products.Product, int, error) {
	products, total, err := ps.productRepo.GetByCategory(category, page, search, sort)
	if err != nil {
		return nil, 0, err
	}
	if page > total {
		return nil, 0, constant.ErrPageInvalid
	}

	return products, total, nil
}

func (ps *ProductService) GetById(id string) (products.Product, error) {
	return ps.productRepo.GetById(id)
}

func (ps *ProductService) Update(product products.Product) error {
	for i, impact := range product.ImpactCategories {
		data, _ := ps.impactRepo.GetByID(impact.ImpactCategoryID)
		if data.ID == "" {
			return constant.ErrCreateProduct
		}
		impact.ID = uuid.New().String()
		product.ImpactCategories[i] = impact
	}

	for i, image := range product.Images {
		image.ID = uuid.New().String()
		product.Images[i] = image
	}

	return ps.productRepo.Update(product)
}

func (ps *ProductService) Delete(productId string) error {
	return ps.productRepo.Delete(productId)
}
