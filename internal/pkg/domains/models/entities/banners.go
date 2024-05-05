package entities

type Banner struct {
	ID    int    `gorm:"column:id;primaryKey;type:bigint;not null;autoIncrement" mapstructure:"id" json:"id"`
	Name  string `json:"name"`
	Image string `json:"image"`
	BaseEntity
}
