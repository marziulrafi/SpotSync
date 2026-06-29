package service

import (
	"errors"
	"spotsync/dto"
	"spotsync/models"
	"spotsync/repository"

	"gorm.io/gorm"
)

var ErrZoneNotFound = errors.New("parking zone not found")

type ZoneService interface {
	Create(req dto.CreateZoneRequest) (*dto.ZoneResponse, error)
	GetAll() ([]dto.ZoneResponse, error)
	GetByID(id uint) (*dto.ZoneResponse, error)
	Update(id uint, req dto.UpdateZoneRequest) (*dto.ZoneResponse, error)
	Delete(id uint) error
}

type zoneService struct {
	zoneRepo repository.ZoneRepository
}

func NewZoneService(zoneRepo repository.ZoneRepository) ZoneService {
	return &zoneService{zoneRepo}
}

func (s *zoneService) Create(req dto.CreateZoneRequest) (*dto.ZoneResponse, error) {
	zone := &models.ParkingZone{
		Name:          req.Name,
		Type:          req.Type,
		TotalCapacity: req.TotalCapacity,
		PricePerHour:  req.PricePerHour,
	}
	if err := s.zoneRepo.Create(zone); err != nil {
		return nil, err
	}
	return s.toResponse(zone, 0), nil
}

func (s *zoneService) GetAll() ([]dto.ZoneResponse, error) {
	zones, err := s.zoneRepo.FindAll()
	if err != nil {
		return nil, err
	}

	result := make([]dto.ZoneResponse, 0, len(zones))
	for _, z := range zones {
		activeCount, err := s.zoneRepo.CountActiveReservations(z.ID)
		if err != nil {
			return nil, err
		}
		available := z.TotalCapacity - int(activeCount)
		result = append(result, *s.toResponse(&z, available))
	}
	return result, nil
}

func (s *zoneService) GetByID(id uint) (*dto.ZoneResponse, error) {
	zone, err := s.zoneRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrZoneNotFound
		}
		return nil, err
	}

	activeCount, err := s.zoneRepo.CountActiveReservations(zone.ID)
	if err != nil {
		return nil, err
	}
	available := zone.TotalCapacity - int(activeCount)
	return s.toResponse(zone, available), nil
}

func (s *zoneService) Update(id uint, req dto.UpdateZoneRequest) (*dto.ZoneResponse, error) {
	zone, err := s.zoneRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrZoneNotFound
		}
		return nil, err
	}

	// Patch only provided fields
	if req.Name != "" {
		zone.Name = req.Name
	}
	if req.Type != "" {
		zone.Type = req.Type
	}
	if req.TotalCapacity > 0 {
		zone.TotalCapacity = req.TotalCapacity
	}
	if req.PricePerHour > 0 {
		zone.PricePerHour = req.PricePerHour
	}

	if err := s.zoneRepo.Update(zone); err != nil {
		return nil, err
	}

	activeCount, _ := s.zoneRepo.CountActiveReservations(zone.ID)
	available := zone.TotalCapacity - int(activeCount)
	return s.toResponse(zone, available), nil
}

func (s *zoneService) Delete(id uint) error {
	if _, err := s.zoneRepo.FindByID(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrZoneNotFound
		}
		return err
	}
	return s.zoneRepo.Delete(id)
}

// toResponse maps ParkingZone model to ZoneResponse DTO
func (s *zoneService) toResponse(zone *models.ParkingZone, availableSpots int) *dto.ZoneResponse {
	return &dto.ZoneResponse{
		ID:             zone.ID,
		Name:           zone.Name,
		Type:           zone.Type,
		TotalCapacity:  zone.TotalCapacity,
		AvailableSpots: availableSpots,
		PricePerHour:   zone.PricePerHour,
		CreatedAt:      zone.CreatedAt,
	}
}
