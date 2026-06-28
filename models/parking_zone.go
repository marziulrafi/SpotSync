package models

import "time"

type ParkingZone struct {
	ID             uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name           string    `gorm:"not null" json:"name"`
	Type           string    `gorm:"not null" json:"type"`
	TotalCapacity  int       `gorm:"not null" json:"total_capacity"`
	PricePerHour   float64   `gorm:"not null;type:decimal(10,2)" json:"price_per_hour"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
