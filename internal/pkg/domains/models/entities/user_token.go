package entities

type UserToken struct {
	ID      int    `gorm:"column:id;primaryKey;type:bigint;not null;autoIncrement" mapstructure:"id" json:"id"`
	UserID  int    `json:"user_id"`
	TokenID string `json:"token_id"`
	BaseEntity
}
