package handler

import (
	"net/http"
	"spotsync/dto"
	"spotsync/service"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	authService service.AuthService
	validate    *validator.Validate
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService, validator.New()}
}

func (h *AuthHandler) Register(c echo.Context) error {
	var req dto.RegisterRequest
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

	user, err := h.authService.Register(req)
	if err != nil {
		if err == service.ErrEmailTaken {
			return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Success: false, Message: "Email already in use",
			})
		}
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false, Message: "Registration failed", Errors: err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, dto.SuccessResponse{
		Success: true, Message: "User registered successfully", Data: user,
	})
}

func (h *AuthHandler) Login(c echo.Context) error {
	var req dto.LoginRequest
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

	data, err := h.authService.Login(req)
	if err != nil {
		if err == service.ErrInvalidCredentials {
			return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
				Success: false, Message: "Invalid email or password",
			})
		}
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false, Message: "Login failed", Errors: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true, Message: "Login successful", Data: data,
	})
}
