package jobs

import (
	"go-server/internal/pkg/domains/models/entities"

	"gorm.io/gorm"
)

func CalcPlaceRateJob(db *gorm.DB) {
	var places []entities.Place
	err := db.Find(&places).Error
	if err != nil {
		return
	}
}
