package entities

type PlaceCategory struct {
	ID         int `gorm:"column:id;primaryKey;type:bigint;not null;autoIncrement" mapstructure:"id" json:"id"`
	PlaceID    int `gorm:"primaryKey;not null" json:"place_id"`
	CategoryID int `gorm:"primaryKey;not null" json:"category_id"`
	BaseEntity
}
