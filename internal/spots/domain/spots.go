package domain

import (
	"github.com/google/uuid"
)

// Spot is a structure to represent the spot user wants
type Spot struct {
	ID          uuid.UUID `json:"spot_id" db:"id"`
	Name        string    `json:"spot_name" db:"name"`
	Website     string    `json:"website" db:"website"`
	Coordinates Geography `json:"coordinates" db:"coordinates"`
	Description string    `json:"description" db:"description"`
	Rating      float32   `json:"rating" db:"rating"`
}

// Geography struct a location, consists of lat and lon
type Geography struct {
	Latitude  float64 `json:"latitude" db:"latitude" validate:"required"`
	Longitude float64 `json:"longitude" db:"longitude" validate:"required"`
}

// SpotList response
type SpotList struct {
	TotalCount int     `json:"total_count"`
	TotalPages int     `json:"total_pages"`
	Page       int     `json:"page"`
	Size       int     `json:"size"`
	HasMore    bool    `json:"has_more"`
	Spots      []*Spot `json:"spots"`
}
