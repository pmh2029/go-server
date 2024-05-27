package entities

import (
	"time"

	"gorm.io/gorm"
)

type Comment struct {
	ID      int    `json:"id" binding:"required"`
	Rate    int    `json:"rate"`
	Comment string `json:"comment"`
	UserID  int    `json:"user_id"`
	User    User   `gorm:"foreignKey:UserID;references:ID" json:"user"`
	PlaceID int    `json:"place_id"`
	Place   Place  `gorm:"foreignKey:PlaceID;references:ID" json:"place"`
	BaseEntity
}

func (i *Comment) AfterFind(tx *gorm.DB) (err error) {
	local, err := time.LoadLocation("Asia/Ho_Chi_Minh")
	if err != nil {
		return
	}

	createdAtTimeLocal, err := time.ParseInLocation("2006-01-02 15:04:05", i.CreatedAt.Format("2006-01-02 15:04:05"), local)
	if err != nil {
		return
	}

	i.CreatedAt = &createdAtTimeLocal

	return
}
