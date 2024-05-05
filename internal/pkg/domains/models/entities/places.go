package entities

type Place struct {
	ID          int        `gorm:"column:id;primaryKey;type:bigint;not null;autoIncrement" mapstructure:"id" json:"id"`
	Name        string     `gorm:"not null" json:"name"`
	Address     string     `gorm:"not null" json:"address"`
	Latitude    float64    `gorm:"not null" json:"latitude"`
	Longitude   float64    `gorm:"not null" json:"longitude"`
	Description string     `json:"description"`
	Images      string     `gorm:"not null" json:"images"`
	Price       float64    `json:"price"`
	Rate        float64    `json:"rate"`
	Categories  []Category `gorm:"many2many:place_categories"`
	BaseEntity
}
