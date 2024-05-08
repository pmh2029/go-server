package dtos

type CreatePlaceRequestDto struct {
	Name        string   `json:"name" binding:"required,min=1"`
	Address     string   `json:"address" binding:"required,min=1"`
	Latitude    float64  `json:"latitude" binding:"required,min=1"`
	Longitude   float64  `json:"longitude" binding:"required,min=1"`
	Description string   `json:"description"`
	Images      []string `json:"images" binding:"required"`
	Price       float64  `json:"price" binding:"required,min=1"`
	Categories  []int    `json:"categories" binding:"required"`
}

type UpdatePlaceRequestDto struct {
	Name        string   `json:"name" binding:"required,min=1"`
	Address     string   `json:"address" binding:"required,min=1"`
	Latitude    float64  `json:"latitude" binding:"required,min=1"`
	Longitude   float64  `json:"longitude" binding:"required,min=1"`
	Description string   `json:"description"`
	Images      []string `json:"images" binding:"required"`
	Price       float64  `json:"price" binding:"required,min=1"`
	Categories  []int    `json:"categories" binding:"required"`
}
