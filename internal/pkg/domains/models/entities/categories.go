package entities

type Category struct {
	ID          int     `gorm:"column:id;primaryKey;type:bigint;not null;autoIncrement" mapstructure:"id" json:"id"`
	Name        string  `gorm:"not null" json:"name,omitempty"`
	Description string  `json:"description,omitempty"`
	Icon        string  `json:"icon,omitempty"`
	Places      []Place `gorm:"many2many:place_categories" json:"places,omitempty"`
	BaseEntity
}
