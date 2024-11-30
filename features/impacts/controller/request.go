package controller

type CreateImpactRequest struct {
	Name        string `json:"name" validate:"required,min=1"`
	ImpactPoint int    `json:"impact_point" validate:"required,min=0"`
	Description string `json:"description" validate:"required,min=1"`
}
