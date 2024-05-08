package database

import (
	"context"

	"gorm.io/gorm"
)

func Transaction(ctx context.Context, db *gorm.DB, callback func(db *gorm.DB) error) error {
	tx := db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}
	committed := false
	defer (func() {
		if !committed {
			tx.Rollback()
		}
	})()
	if err := callback(tx); err != nil {
		return err
	}
	if err := tx.Commit().Error; err != nil {
		return err
	}
	committed = true
	return nil
}
