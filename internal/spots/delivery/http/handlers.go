package http

import (
	"net/http"

	"github.com/kkr2/spots/internal/config"
	"github.com/kkr2/spots/internal/spots/service"
	"github.com/kkr2/spots/pkg/httpErrors"
	"github.com/kkr2/spots/pkg/logger"
	"github.com/kkr2/spots/pkg/utils"
	"github.com/labstack/echo/v4"
)
// SpotsHandlers is an http interface 
type SpotsHandlers interface {
	GetSpots() echo.HandlerFunc
}

// spotsHandlers is the concrete implementation of spots handler 
type spotsHandlers struct {
	cfg         *config.Config
	spotService service.SpotService
	logger      logger.Logger
}

// NewSpotsHandlers News handlers constructor
func NewSpotsHandlers(cfg *config.Config, spotService service.SpotService, logger logger.Logger) SpotsHandlers {
	return &spotsHandlers{cfg: cfg, spotService: spotService, logger: logger}
}

// GetSpots is a handler function using GetSpotsInRangeOfCoordinate use-case
func (h spotsHandlers) GetSpots() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := utils.GetRequestCtx(c)
		pq, err := utils.GetPaginationFromCtx(c)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}
		coordinate, err := utils.GetCoordinateFromCtx(c)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}
		spotList, err := h.spotService.GetSpotsInRangeOfCoordinate(ctx, coordinate, pq)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, spotList)
	}
}
