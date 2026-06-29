package handler

import (
	"net/http"
	"strconv"
	"spotsync/dto"
	"spotsync/service"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type ZoneHandler struct {
	zoneService service.ZoneService
	validate    *validator.Validate
}

func NewZoneHandler(zoneService service.ZoneService) *ZoneHandler {
	return &ZoneHandler{zoneService, validator.New()}
}

func (h *ZoneHandler) Create(c echo.Context) error {
	var req dto.CreateZoneRequest
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

	zone, err := h.zoneService.Create(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false, Message: "Failed to create parking zone", Errors: err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, dto.SuccessResponse{
		Success: true, Message: "Parking zone created successfully", Data: zone,
	})
}

func (h *ZoneHandler) GetAll(c echo.Context) error {
	zones, err := h.zoneService.GetAll()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false, Message: "Failed to retrieve parking zones", Errors: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true, Message: "Parking zones retrieved successfully", Data: zones,
	})
}

func (h *ZoneHandler) GetByID(c echo.Context) error {
	id, err := parseID(c, "id")
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false, Message: "Invalid zone ID",
		})
	}

	zone, err := h.zoneService.GetByID(uint(id))
	if err != nil {
		if err == service.ErrZoneNotFound {
			return c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Success: false, Message: "Parking zone not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false, Message: "Failed to retrieve zone", Errors: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true, Message: "Parking zone retrieved successfully", Data: zone,
	})
}

func (h *ZoneHandler) Update(c echo.Context) error {
	id, err := parseID(c, "id")
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false, Message: "Invalid zone ID",
		})
	}

	var req dto.UpdateZoneRequest
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

	zone, err := h.zoneService.Update(uint(id), req)
	if err != nil {
		if err == service.ErrZoneNotFound {
			return c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Success: false, Message: "Parking zone not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false, Message: "Failed to update zone", Errors: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true, Message: "Parking zone updated successfully", Data: zone,
	})
}

func (h *ZoneHandler) Delete(c echo.Context) error {
	id, err := parseID(c, "id")
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false, Message: "Invalid zone ID",
		})
	}

	if err := h.zoneService.Delete(uint(id)); err != nil {
		if err == service.ErrZoneNotFound {
			return c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Success: false, Message: "Parking zone not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false, Message: "Failed to delete zone", Errors: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true, Message: "Parking zone deleted successfully",
	})
}

// parseID is a shared helper to extract and validate uint path params
func parseID(c echo.Context, param string) (uint64, error) {
	return strconv.ParseUint(c.Param(param), 10, 64)
}
