package entities

type Place struct {
	ID          int
	Name        string
	Address     string
	Latitude    float64
	Longitude   float64
	Description string
	Image       string
	Price       float64
	Rate        float64
	Categories  []Category `gorm:"many2many:place_categories"`
	BaseEntity
}
