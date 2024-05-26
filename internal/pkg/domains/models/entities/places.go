package entities

import (
	"strings"

	"gorm.io/gorm"
)

type Place struct {
	ID             int        `gorm:"column:id;primaryKey;type:bigint;not null;autoIncrement" mapstructure:"id" json:"id"`
	Name           string     `gorm:"not null" json:"name,omitempty"`
	Address        string     `gorm:"not null" json:"address,omitempty"`
	Latitude       float64    `gorm:"not null" json:"latitude,omitempty"`
	Longitude      float64    `gorm:"not null" json:"longitude,omitempty"`
	Description    string     `json:"description,omitempty"`
	Images         string     `gorm:"not null" json:"-"`
	Price          float64    `json:"price,omitempty"`
	Rate           float64    `json:"rate,omitempty"`
	Categories     []Category `gorm:"many2many:place_categories" json:"categories,omitempty"`
	ImagesResponse []string   `gorm:"-" json:"images,omitempty"`
	BaseEntity
}

func (i *Place) AfterFind(tx *gorm.DB) (err error) {
	i.ImagesResponse = strings.Split(i.Images, "|")
	i.ImagesResponse = i.ImagesResponse[1 : len(i.ImagesResponse)-1]

	var comment []Comment
	err = tx.Where("place_id = ?", i.ID).Find(&comment).Error
	if err != nil {
		return
	}
	if len(comment) == 0 {
		i.Rate = 0
	} else {
		for _, v := range comment {
			i.Rate += float64(v.Rate)
		}
		i.Rate /= float64(len(comment))
	}
	return
}
