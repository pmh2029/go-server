package dtos

import "go-server/internal/pkg/domains/models/entities"

type CreateBannerRequestDto struct {
	Name  string `json:"name" binding:"required"`
	Image string `json:"image" binding:"required"`
}

type CreateBannerResponseDto struct {
	Banner entities.Banner `json:"banner"`
}

type UpdateBannerRequestDto struct {
	Name  string `json:"name"`
	Image string `json:"image"`
}

type UpdateBannerResponseDto struct {
	Banner entities.Banner `json:"banner"`
}
