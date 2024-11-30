package service

import (
	"greenenvironment/constant"
	"greenenvironment/features/impacts"
)

type ImpactService struct {
	impactRepo impacts.ImpactRepositoryInterface
}

func NewNewImpactService(d impacts.ImpactRepositoryInterface) impacts.ImpactServiceInterface {
	return &ImpactService{
		impactRepo: d,
	}
}

func (is *ImpactService) GetAll() ([]impacts.ImpactCategory, error) {
	return is.impactRepo.GetAll()
}

func (is *ImpactService) GetByID(ID string) (impacts.ImpactCategory, error) {
	if ID == "" {
		return impacts.ImpactCategory{}, constant.ErrImpactCategoryNotFound
	}
	return is.impactRepo.GetByID(ID)
}

func (is *ImpactService) Create(category impacts.ImpactCategory) error {

	return is.impactRepo.Create(category)
}

func (is *ImpactService) Delete(category impacts.ImpactCategory) error {
	return is.impactRepo.Delete(category)
}
