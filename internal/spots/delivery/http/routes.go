package http

import "github.com/labstack/echo/v4"

// MapSpotsRoutes maps spot endpoints to group
func MapSpotsRoutes(spotGroup *echo.Group, h SpotsHandlers) {
	spotGroup.GET("", h.GetSpots())
}
