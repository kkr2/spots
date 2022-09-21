package utils

import (
	"strconv"

	"github.com/kkr2/spots/internal/spots/domain"
	"github.com/labstack/echo/v4"
)

func GetCoordinateFromCtx(c echo.Context) (*domain.Geography, error) {

	lat, _ := strconv.ParseFloat(c.QueryParam("lat"), 64)

	lon, _ := strconv.ParseFloat(c.QueryParam("lon"), 64)

	geoPoint := domain.Geography{
		Latitude:  lat,
		Longitude: lon,
	}

	return &geoPoint, nil
}
