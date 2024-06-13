package dtos

import "go-server/internal/pkg/domains/models/entities"

type CreateTripRequestDto struct {
	Owner    int                   `json:"owner" binding:"required,min=1"`
	Name     string                `json:"name" binding:"required,min=1"`
	FromDate int                   `json:"from_date" binding:"required,min=1"`
	ToDate   int                   `json:"to_date" binding:"required,min=1"`
	Users    int                   `json:"users"`
	Days     []CreateDayRequestDto `json:"days" binding:"required"`
}

type CreateDayRequestDto struct {
	Places []entities.PlaceInDay `json:"places"`
}

type UpdateTripRequestDto struct {
	Name     string                `json:"name" binding:"required,min=1"`
	FromDate int                   `json:"from_date" binding:"required,min=1"`
	ToDate   int                   `json:"to_date" binding:"required,min=1"`
	Users    int                   `json:"users"`
	Days     []CreateDayRequestDto `json:"days" binding:"required"`
}
