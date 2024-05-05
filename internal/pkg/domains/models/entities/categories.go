package entities

type Category struct {
	ID          int     `gorm:"column:id;primaryKey;type:bigint;not null;autoIncrement" mapstructure:"id" json:"id"`
	Name        string  `gorm:"not null" json:"name"`
	Description string  `json:"description"`
	Icon        string  `json:"icon"`
	Places      []Place `gorm:"many2many:place_categories"`
	BaseEntity
}
