package repository

import (
	"greenenvironment/constant"
	"greenenvironment/features/impacts"

	"gorm.io/gorm"
)

type ImpactRepository struct {
	DB *gorm.DB
}

func NewImpactRepository(db *gorm.DB) impacts.ImpactRepositoryInterface {
	return &ImpactRepository{
		DB: db,
	}
}
func (ir *ImpactRepository) GetAll() ([]impacts.ImpactCategory, error) {
	var categories []impacts.ImpactCategory
	err := ir.DB.Model(&ImpactCategory{}).Find(&categories).Error
	if err != nil {
		return nil, constant.ErrImpactCategoryNotFound
	}
	return categories, nil
}
func (ir *ImpactRepository) GetByID(ID string) (impacts.ImpactCategory, error) {
	var category impacts.ImpactCategory
	err := ir.DB.First(&category, "id = ?", ID).Error
	if err != nil {
		return impacts.ImpactCategory{}, constant.ErrImpactCategoryNotFound
	}
	return category, nil
}
func (ir *ImpactRepository) Create(category impacts.ImpactCategory) error {
	impactData := ImpactCategory{
		ID:          category.ID,
		Name:        category.Name,
		ImpactPoint: category.ImpactPoint,
		Description: category.Description,
	}
	err := ir.DB.Create(&impactData).Error
	if err != nil {
		return constant.ErrCreateImpactCategory
	}
	return nil
}
func (ir *ImpactRepository) Delete(category impacts.ImpactCategory) error {
	result := ir.DB.Where("id = ?", category.ID).Delete(&ImpactCategory{})
	if result.Error != nil {
		return constant.ErrDeleteImpactCategory
	}
	if result.RowsAffected == 0 {
		return constant.ErrImpactCategoryNotFound
	}
	return nil
}
