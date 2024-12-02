package repository

import (
	"greenenvironment/features/guest"
	"greenenvironment/features/products"
	productRepo "greenenvironment/features/products/repository"

	"gorm.io/gorm"
)

type GuestRepository struct {
	DB *gorm.DB
}

func NewGuestRepository(db *gorm.DB) guest.GuestRepostoryInterface {
	return &GuestRepository{DB: db}
}
func (gr *GuestRepository) GetGuests() (guest.Guest, error) {
	totalProduct, err := gr.GetTotalProduct()
	if err != nil {
		return guest.Guest{}, err
	}
	totalNewProductThisMonth, err := gr.GetTotalNewProductThisMonth()
	if err != nil {
		return guest.Guest{}, err
	}

	newProduct, err := gr.GetNewProduct()
	if err != nil {
		return guest.Guest{}, err
	}

	return guest.Guest{
		TotalProduct:             totalProduct,
		TotalNewProductThisMonth: totalNewProductThisMonth,
		NewProduct:               newProduct,
	}, nil
}
func (gr *GuestRepository) GetTotalProduct() (int, error) {
	var totalProduct int64
	err := gr.DB.Table("products").Count(&totalProduct).Error
	if err != nil {
		return 0, err
	}
	return int(totalProduct), nil

}
func (gr *GuestRepository) GetTotalNewProductThisMonth() (int, error) {
	var totalNewProductThisMonth int64
	err := gr.DB.Table("products").
		Where("MONTH(created_at) = MONTH(NOW())").
		Count(&totalNewProductThisMonth).Error
	if err != nil {
		return 0, err
	}
	return int(totalNewProductThisMonth), nil
}
func (gr *GuestRepository) GetNewProduct() ([]products.Product, error) {
	var top10Products []products.Product
	err := gr.DB.Model(&productRepo.Product{}).
		Preload("Images").
		Preload("ImpactCategories.ImpactCategory").
		Order("created_at DESC").
		Limit(10).
		Find(&top10Products).Error
	if err != nil {
		return nil, err
	}

	return top10Products, nil
}
