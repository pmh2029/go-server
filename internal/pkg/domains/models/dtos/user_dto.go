package dtos

import "go-server/internal/pkg/domains/models/entities"

type RegisterRequestDto struct {
	Username string `json:"username" binding:"required,alphaNumeric"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterResponseDto struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type LoginRequestDto struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponseDto struct {
	User        entities.User `json:"user"`
	AccessToken string        `json:"access_token"`
}
