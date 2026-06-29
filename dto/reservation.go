package dto

import "time"

type CreateReservationRequest struct {
	ZoneID       uint   `json:"zone_id"       validate:"required"`
	LicensePlate string `json:"license_plate" validate:"required,max=15"`
}

type MyReservationResponse struct {
	ID           uint      `json:"id"`
	LicensePlate string    `json:"license_plate"`
	Status       string    `json:"status"`
	Zone         ZoneInfo  `json:"zone"`
	CreatedAt    time.Time `json:"created_at"`
}

type ZoneInfo struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type ReservationResponse struct {
	ID           uint      `json:"id"`
	UserID       uint      `json:"user_id"`
	ZoneID       uint      `json:"zone_id"`
	LicensePlate string    `json:"license_plate"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type AdminReservationResponse struct {
	ID           uint         `json:"id"`
	LicensePlate string       `json:"license_plate"`
	Status       string       `json:"status"`
	User         UserInfo     `json:"user"`
	Zone         ZoneInfo     `json:"zone"`
	CreatedAt    time.Time    `json:"created_at"`
}

type UserInfo struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
