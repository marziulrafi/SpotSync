package repository

import (
	"spotsync/models"

	"gorm.io/gorm"
)

type ZoneRepository interface {
	Create(zone *models.ParkingZone) error
	FindAll() ([]models.ParkingZone, error)
	FindByID(id uint) (*models.ParkingZone, error)
	Update(zone *models.ParkingZone) error
	Delete(id uint) error
	CountActiveReservations(zoneID uint) (int64, error)
}

type zoneRepository struct {
	db *gorm.DB
}

func NewZoneRepository(db *gorm.DB) ZoneRepository {
	return &zoneRepository{db}
}

func (r *zoneRepository) Create(zone *models.ParkingZone) error {
	return r.db.Create(zone).Error
}

func (r *zoneRepository) FindAll() ([]models.ParkingZone, error) {
	var zones []models.ParkingZone
	if err := r.db.Find(&zones).Error; err != nil {
		return nil, err
	}
	return zones, nil
}

func (r *zoneRepository) FindByID(id uint) (*models.ParkingZone, error) {
	var zone models.ParkingZone
	if err := r.db.First(&zone, id).Error; err != nil {
		return nil, err
	}
	return &zone, nil
}

func (r *zoneRepository) Update(zone *models.ParkingZone) error {
	return r.db.Save(zone).Error
}

func (r *zoneRepository) Delete(id uint) error {
	return r.db.Delete(&models.ParkingZone{}, id).Error
}

// CountActiveReservations returns the count of active reservations for a zone
func (r *zoneRepository) CountActiveReservations(zoneID uint) (int64, error) {
	var count int64
	if err := r.db.Model(&models.Reservation{}).
		Where("zone_id = ? AND status = 'active'", zoneID).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
