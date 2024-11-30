package controller

type ImpactCategoryResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	ImpactPoint int    `json:"impact_point"`
	Description string `json:"description"`
}
