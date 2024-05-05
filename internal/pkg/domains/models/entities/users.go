package entities

import (
	"time"
)

// UsersTableName TableName
var UsersTableName = "users"

type User struct {
	ID       int        `gorm:"column:id;primaryKey;type:bigint;not null;autoIncrement" mapstructure:"id" json:"id"`
	Username string     `gorm:"column:username;not null;unique" json:"username,omitempty"`
	Email    string     `gorm:"column:email;not null;unique" json:"email,omitempty"`
	Active   bool       `json:"active,omitempty"`
	Avatar   string     `json:"avatar,omitempty"`
	BirthDay *time.Time `json:"birth_day,omitempty"`
	Gender   int        `json:"gender,omitempty"` // 1: nam, 2: nữ
	Contact  string     `json:"contact,omitempty"`
	Password string     `gorm:"column:password;not null" json:"password,omitempty"`
	IsAdmin  bool       `gorm:"default:false"`
	BaseEntity
}

// TableName func
func (i *User) TableName() string {
	return UsersTableName
}
