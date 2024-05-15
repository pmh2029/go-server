package entities

type UserTrip struct {
	TripID int `gorm:"primaryKey"`
	UserID int `gorm:"primaryKey"`
}
