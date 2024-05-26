package dtos

import "go-server/internal/pkg/domains/models/entities"

type RegisterRequestDto struct {
	Username string `json:"username" binding:"required,alphaNumeric"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterResponseDto struct {
	User        entities.User `json:"user"`
	AccessToken string        `json:"access_token"`
}

type LoginRequestDto struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponseDto struct {
	User        entities.User `json:"user"`
	AccessToken string        `json:"access_token"`
}

type UpdateUserRequestDto struct {
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
	BirthDay int    `json:"birth_day"`
	Gender   int    `json:"gender"`
	Contact  string `json:"contact"`
}

type UpdateUserResponseDto struct {
	User entities.User `json:"user"`
}

type ForgotPasswordRequestDto struct {
	Email        string `json:"email" binding:"required,email"`
	ReceiveEmail string `json:"receive_email" binding:"required,email"`
}

type ChangePasswordRequestDto struct {
	OldPassword     string `json:"old_password" binding:"required,min=1"`
	NewPassword     string `json:"new_password" binding:"required,min=1"`
	ConfirmPassword string `json:"confirm_password" binding:"required,min=1"`
}

type UpdateStatusRequestDto struct {
	UserID int `json:"user_id" binding:"required,min=1"`
	Status int `json:"status" binding:"required,min=1,max=2"`
}
