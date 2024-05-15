package entities

import (
	"encoding/json"

	"gorm.io/gorm"
)

type PlaceInDay struct {
	ID        int    `json:"id" binding:"required"`
	Note      string `json:"note"`
	VisitTime int    `json:"visit_time"`
	StartTime int    `json:"start_time"`
	Vehicle   int    `json:"vehicle"`
}

type Day struct {
	ID         int          `gorm:"column:id;primaryKey;type:bigint;not null;autoIncrement" mapstructure:"id" json:"id"`
	Places     string       `json:"-"`
	TripID     int          `json:"trip_id"`
	PlacesJson []PlaceInDay `gorm:"-" json:"places"`
	BaseEntity
}

func (i *Day) AfterFind(tx *gorm.DB) (err error) {
	err = json.Unmarshal([]byte(i.Places), &i.PlacesJson)

	return
}
