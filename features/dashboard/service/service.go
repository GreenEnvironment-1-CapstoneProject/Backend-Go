package service

import "greenenvironment/features/dashboard"

type DashboardService struct {
	repo dashboard.DashboardRepositoryInterface
}

func NewDashboardService(repo dashboard.DashboardRepositoryInterface) dashboard.DashboardServiceInterface {
	return &DashboardService{repo: repo}
}

func (ds *DashboardService) GetDashboardData(filter string) (dashboard.DashboardData, error) {
	data, err := ds.repo.GetDashboardData(filter)
	if err != nil {
		return dashboard.DashboardData{}, err
	}

	topCategories, err := ds.repo.GetTopCategories(filter)
	if err != nil {
		return dashboard.DashboardData{}, err
	}
	data.TopCategories = topCategories

	lastTransactions, err := ds.repo.GetLastTransactions()
	if err != nil {
		return dashboard.DashboardData{}, err
	}
	data.LastTransactions = lastTransactions

	return data, nil
}

