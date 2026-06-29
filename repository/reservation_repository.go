package repository

import (
	"errors"
	"spotsync/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ErrZoneFull is returned when a zone has no available spots
var ErrZoneFull = errors.New("parking zone is at full capacity")

type ReservationRepository interface {
	// CreateWithLock uses a DB transaction + FOR UPDATE row lock to safely check capacity
	CreateWithLock(reservation *models.Reservation) error
	FindByUserID(userID uint) ([]models.Reservation, error)
	FindByID(id uint) (*models.Reservation, error)
	Cancel(id uint) error
	FindAll() ([]models.Reservation, error)
}

type reservationRepository struct {
	db *gorm.DB
}

func NewReservationRepository(db *gorm.DB) ReservationRepository {
	return &reservationRepository{db}
}

// CreateWithLock atomically checks capacity and creates the reservation.
// This is the solution to the "EV Spot Bottleneck" race condition.
func (r *reservationRepository) CreateWithLock(reservation *models.Reservation) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Step 1: Acquire a row-level lock on the zone record.
		// "FOR UPDATE" prevents any other transaction from reading or
		// modifying this row until we commit — eliminating the race condition.
		var zone models.ParkingZone
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&zone, reservation.ZoneID).Error; err != nil {
			return err
		}

		// Step 2: Count current active reservations for this zone
		var activeCount int64
		if err := tx.Model(&models.Reservation{}).
			Where("zone_id = ? AND status = 'active'", reservation.ZoneID).
			Count(&activeCount).Error; err != nil {
			return err
		}

		// Step 3: Enforce capacity limit
		if activeCount >= int64(zone.TotalCapacity) {
			return ErrZoneFull
		}

		// Step 4: Safe to create — still within capacity
		return tx.Create(reservation).Error
	})
}

func (r *reservationRepository) FindByUserID(userID uint) ([]models.Reservation, error) {
	var reservations []models.Reservation
	if err := r.db.Where("user_id = ?", userID).
		Preload("Zone").
		Order("created_at DESC").
		Find(&reservations).Error; err != nil {
		return nil, err
	}
	return reservations, nil
}

func (r *reservationRepository) FindByID(id uint) (*models.Reservation, error) {
	var res models.Reservation
	if err := r.db.First(&res, id).Error; err != nil {
		return nil, err
	}
	return &res, nil
}

// Cancel sets the reservation status to "cancelled" (soft-cancel, not hard delete)
func (r *reservationRepository) Cancel(id uint) error {
	return r.db.Model(&models.Reservation{}).
		Where("id = ?", id).
		Update("status", "cancelled").Error
}

func (r *reservationRepository) FindAll() ([]models.Reservation, error) {
	var reservations []models.Reservation
	if err := r.db.
		Preload("User").
		Preload("Zone").
		Order("created_at DESC").
		Find(&reservations).Error; err != nil {
		return nil, err
	}
	return reservations, nil
}
