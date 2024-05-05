package dtos

type CreateCategoryRequestDto struct {
	Name        string `json:"name" binding:"required"`
	Icon        string `json:"icon" binding:"required"`
	Description string `json:"description"`
}

type UpdateCategoryRequestDto struct {
	Name        string `json:"name"`
	Icon        string `json:"icon"`
	Description string `json:"description"`
}
