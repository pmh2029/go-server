package dtos

type CreateCommentRequestDto struct {
	Rate    int `json:"rate"`
	Comment string  `json:"comment"`
	PlaceID int     `json:"place_id" binding:"required"`
	UserID  int     `json:"user_id" binding:"required"`
}

type UpdateCommentRequestDto struct {
	Rate    int `json:"rate"`
	Comment string  `json:"comment"`
}
