package entities

type Comment struct {
	ID      int     `json:"id" binding:"required"`
	Rate    float64 `json:"rate"`
	Comment string  `json:"comment"`
	UserID  int     `json:"user_id"`
	User    User    `gorm:"foreignKey:UserID;references:ID" json:"user"`
	PlaceID int     `json:"place_id"`
	Place   Place   `gorm:"foreignKey:PlaceID;references:ID" json:"place"`
	BaseEntity
}
