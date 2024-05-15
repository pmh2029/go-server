package entities

import (
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

type Trip struct {
	ID           int       `gorm:"column:id;primaryKey;type:bigint;not null;autoIncrement" mapstructure:"id" json:"id"`
	Days         []Day     `gorm:"foreignKey:TripID" json:"days"`
	FromDate     time.Time `json:"-"`
	Name         string    `json:"name"`
	ToDate       time.Time `json:"-"`
	Users        []User    `gorm:"many2many:user_trips" json:"users"`
	Owner        int       `json:"owner"`
	FromDateUnix int       `gorm:"-" json:"from_date"`
	UserIDs      string    `json:"-"`
	ToDateUnix   int       `gorm:"-" json:"to_date"`
	BaseEntity
}

func (i *Trip) AfterFind(tx *gorm.DB) (err error) {
	i.FromDateUnix = int(i.FromDate.Unix())
	i.ToDateUnix = int(i.ToDate.Unix())

	if i.UserIDs != "" {
		userIDs := strings.Split(i.UserIDs, ",")
		userIDs = userIDs[1 : len(userIDs)-1]
		userOrder := strings.Join(userIDs, ",")

		var users []User
		err = tx.Where("id IN (?)", userIDs).Order(fmt.Sprintf("FIELD(id, %s)", userOrder)).Find(&users).Error
		if err != nil {
			return
		}

		i.Users = users
	}
	return
}
