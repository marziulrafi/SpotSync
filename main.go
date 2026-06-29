package main

import (
	"log"
	"net/http"
	"os"

	"spotsync/config"
	"spotsync/handler"
	mw "spotsync/middleware"
	"spotsync/repository"
	"spotsync/service"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, reading from environment")
	}

	// Connect & auto-migrate database
	db := config.ConnectDB()

	// Repositories
	userRepo := repository.NewUserRepository(db)
	zoneRepo := repository.NewZoneRepository(db)
	resRepo := repository.NewReservationRepository(db)

	// Services
	authSvc := service.NewAuthService(userRepo)
	zoneSvc := service.NewZoneService(zoneRepo)
	resSvc := service.NewReservationService(resRepo, zoneRepo)

	// Handlers
	authH := handler.NewAuthHandler(authSvc)
	zoneH := handler.NewZoneHandler(zoneSvc)
	resH := handler.NewReservationHandler(resSvc)

	e := echo.New()

	// Global middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete, http.MethodOptions},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))

	api := e.Group("/api/v1")

	// Auth (public)
	auth := api.Group("/auth")
	auth.POST("/register", authH.Register)
	auth.POST("/login", authH.Login)

	// Zones
	zones := api.Group("/zones")
	zones.GET("", zoneH.GetAll)
	zones.GET("/:id", zoneH.GetByID)
	zones.POST("", zoneH.Create, mw.JWTMiddleware, mw.RequireRole("admin"))
	zones.PUT("/:id", zoneH.Update, mw.JWTMiddleware, mw.RequireRole("admin"))
	zones.DELETE("/:id", zoneH.Delete, mw.JWTMiddleware, mw.RequireRole("admin"))

	// Reservations
	res := api.Group("/reservations", mw.JWTMiddleware)
	res.POST("", resH.Create)
	res.GET("/my-reservations", resH.GetMyReservations)
	res.DELETE("/:id", resH.Cancel)
	res.GET("", resH.GetAll, mw.RequireRole("admin"))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚗 SpotSync running on :%s", port)
	e.Logger.Fatal(e.Start(":" + port))
}
