package entities

type Category struct {
	ID          int `gorm:"column:id;primaryKey;type:bigint;not null;autoIncrement" mapstructure:"id" json:"id"`
	Name        string
	Description string
	Icon        string
	BaseEntity
}
