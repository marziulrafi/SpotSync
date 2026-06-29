package handler

import (
	"net/http"
	"spotsync/dto"
	"spotsync/repository"
	"spotsync/service"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type ReservationHandler struct {
	resService service.ReservationService
	validate   *validator.Validate
}

func NewReservationHandler(resService service.ReservationService) *ReservationHandler {
	return &ReservationHandler{resService, validator.New()}
}

func (h *ReservationHandler) Create(c echo.Context) error {
	userID := c.Get("user_id").(uint)

	var req dto.CreateReservationRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false, Message: "Invalid request body", Errors: err.Error(),
		})
	}
	if err := h.validate.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false, Message: "Validation failed", Errors: err.Error(),
		})
	}

	res, err := h.resService.Create(userID, req)
	if err != nil {
		switch err {
		case service.ErrZoneNotFound:
			return c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Success: false, Message: "Parking zone not found",
			})
		case repository.ErrZoneFull:
			return c.JSON(http.StatusConflict, dto.ErrorResponse{
				Success: false, Message: "Parking zone is at full capacity",
			})
		default:
			return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Success: false, Message: "Failed to create reservation", Errors: err.Error(),
			})
		}
	}

	return c.JSON(http.StatusCreated, dto.SuccessResponse{
		Success: true, Message: "Reservation confirmed successfully", Data: res,
	})
}

func (h *ReservationHandler) GetMyReservations(c echo.Context) error {
	userID := c.Get("user_id").(uint)

	reservations, err := h.resService.GetMyReservations(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false, Message: "Failed to retrieve reservations", Errors: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true, Message: "My reservations retrieved successfully", Data: reservations,
	})
}

func (h *ReservationHandler) Cancel(c echo.Context) error {
	userID := c.Get("user_id").(uint)
	role := c.Get("role").(string)

	id, err := parseID(c, "id")
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false, Message: "Invalid reservation ID",
		})
	}

	if err := h.resService.Cancel(uint(id), userID, role); err != nil {
		switch err {
		case service.ErrReservationNotFound:
			return c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Success: false, Message: "Reservation not found",
			})
		case service.ErrForbidden:
			return c.JSON(http.StatusForbidden, dto.ErrorResponse{
				Success: false, Message: "You can only cancel your own reservations",
			})
		case service.ErrAlreadyCancelled:
			return c.JSON(http.StatusConflict, dto.ErrorResponse{
				Success: false, Message: "Reservation is already cancelled",
			})
		default:
			return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Success: false, Message: "Failed to cancel reservation", Errors: err.Error(),
			})
		}
	}

	return c.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true, Message: "Reservation cancelled successfully",
	})
}

func (h *ReservationHandler) GetAll(c echo.Context) error {
	reservations, err := h.resService.GetAll()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false, Message: "Failed to retrieve reservations", Errors: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true, Message: "All reservations retrieved successfully", Data: reservations,
	})
}
