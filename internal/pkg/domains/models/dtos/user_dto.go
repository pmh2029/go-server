package dtos

type RegisterRequestDto struct {
	Username string `json:"username" binding:"required,alphaNumeric"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterResponseDto struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}
