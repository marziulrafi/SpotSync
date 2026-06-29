package service

import (
	"errors"
	"spotsync/dto"
	"spotsync/models"
	"spotsync/repository"

	"gorm.io/gorm"
)

var (
	ErrReservationNotFound = errors.New("reservation not found")
	ErrForbidden           = errors.New("you do not have permission to perform this action")
	ErrAlreadyCancelled    = errors.New("reservation is already cancelled")
)

type ReservationService interface {
	Create(userID uint, req dto.CreateReservationRequest) (*dto.ReservationResponse, error)
	GetMyReservations(userID uint) ([]dto.MyReservationResponse, error)
	Cancel(reservationID, requesterID uint, requesterRole string) error
	GetAll() ([]dto.AdminReservationResponse, error)
}

type reservationService struct {
	resRepo  repository.ReservationRepository
	zoneRepo repository.ZoneRepository
}

func NewReservationService(
	resRepo repository.ReservationRepository,
	zoneRepo repository.ZoneRepository,
) ReservationService {
	return &reservationService{resRepo, zoneRepo}
}

func (s *reservationService) Create(userID uint, req dto.CreateReservationRequest) (*dto.ReservationResponse, error) {
	// Verify zone exists before attempting lock
	if _, err := s.zoneRepo.FindByID(req.ZoneID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrZoneNotFound
		}
		return nil, err
	}

	res := &models.Reservation{
		UserID:       userID,
		ZoneID:       req.ZoneID,
		LicensePlate: req.LicensePlate,
		Status:       "active",
	}

	// This internally handles the transaction + FOR UPDATE row lock
	if err := s.resRepo.CreateWithLock(res); err != nil {
		return nil, err
	}

	return &dto.ReservationResponse{
		ID:           res.ID,
		UserID:       res.UserID,
		ZoneID:       res.ZoneID,
		LicensePlate: res.LicensePlate,
		Status:       res.Status,
		CreatedAt:    res.CreatedAt,
		UpdatedAt:    res.UpdatedAt,
	}, nil
}

func (s *reservationService) GetMyReservations(userID uint) ([]dto.MyReservationResponse, error) {
	reservations, err := s.resRepo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}

	result := make([]dto.MyReservationResponse, 0, len(reservations))
	for _, r := range reservations {
		result = append(result, dto.MyReservationResponse{
			ID:           r.ID,
			LicensePlate: r.LicensePlate,
			Status:       r.Status,
			Zone: dto.ZoneInfo{
				ID:   r.Zone.ID,
				Name: r.Zone.Name,
				Type: r.Zone.Type,
			},
			CreatedAt: r.CreatedAt,
		})
	}
	return result, nil
}

func (s *reservationService) Cancel(reservationID, requesterID uint, requesterRole string) error {
	res, err := s.resRepo.FindByID(reservationID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrReservationNotFound
		}
		return err
	}

	// Drivers can only cancel their own reservations; admins can cancel any
	if requesterRole != "admin" && res.UserID != requesterID {
		return ErrForbidden
	}

	if res.Status == "cancelled" {
		return ErrAlreadyCancelled
	}

	return s.resRepo.Cancel(reservationID)
}

func (s *reservationService) GetAll() ([]dto.AdminReservationResponse, error) {
	reservations, err := s.resRepo.FindAll()
	if err != nil {
		return nil, err
	}

	result := make([]dto.AdminReservationResponse, 0, len(reservations))
	for _, r := range reservations {
		result = append(result, dto.AdminReservationResponse{
			ID:           r.ID,
			LicensePlate: r.LicensePlate,
			Status:       r.Status,
			User: dto.UserInfo{
				ID:    r.User.ID,
				Name:  r.User.Name,
				Email: r.User.Email,
			},
			Zone: dto.ZoneInfo{
				ID:   r.Zone.ID,
				Name: r.Zone.Name,
				Type: r.Zone.Type,
			},
			CreatedAt: r.CreatedAt,
		})
	}
	return result, nil
}
