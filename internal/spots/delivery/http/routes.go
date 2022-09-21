package http

import "github.com/labstack/echo/v4"



func MapSpotsRoutes(spotGroup *echo.Group, h SpotsHandlers) {
	spotGroup.GET("", h.GetSpots())
}
