package server

import (
	"net/http"

	spotHttp "github.com/kkr2/spots/internal/spots/delivery/http"
	"github.com/kkr2/spots/internal/spots/repository"
	"github.com/kkr2/spots/internal/spots/service"
	"github.com/labstack/echo/v4"
)

// MapHandlers glues all code relted to paths
func (s *Server) MapHandlers(e *echo.Echo) error {

	// Init repositories

	sRepo := repository.NewSpotRepository(s.db)

	// Init useCases
	sService := service.NewSpotsService(s.cfg, sRepo, s.slogger)

	// Init handlers
	sHandler := spotHttp.NewSpotsHandlers(s.cfg, sService, s.slogger)

	v1 := e.Group("/api/v1")

	health := v1.Group("/health")

	spotGroup := v1.Group("/spots")

	spotHttp.MapSpotsRoutes(spotGroup, sHandler)

	health.GET("", func(c echo.Context) error {
		s.slogger.Infof("Health check RequestID: %s", c.Response().Header().Get(echo.HeaderXRequestID))
		return c.JSON(http.StatusOK, map[string]string{"status": "OK"})
	})

	return nil
}
