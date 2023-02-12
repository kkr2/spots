package service

import (
	"context"

	"github.com/kkr2/spots/internal/config"
	"github.com/kkr2/spots/internal/spots/domain"
	"github.com/kkr2/spots/internal/spots/repository"
	"github.com/kkr2/spots/pkg/httpErrors"
	"github.com/kkr2/spots/pkg/logger"
	"github.com/kkr2/spots/pkg/utils"
	"github.com/pkg/errors"
)

type SpotService interface {
	GetSpotsInRangeOfCoordinate(ctx context.Context, coordinate *domain.Geography, pq *utils.PaginationQuery) (*domain.SpotList, error)
}

type spotService struct {
	cfg      *config.Config
	spotRepo repository.SpotRepository
	logger   logger.Logger
}

// NewSpotsService Spots service constructor
func NewSpotsService(cfg *config.Config, spotRepo repository.SpotRepository, logger logger.Logger) SpotService {
	return &spotService{cfg: cfg, spotRepo: spotRepo, logger: logger}
}

func (sps *spotService) GetSpotsInRangeOfCoordinate(
	ctx context.Context,
	coordinate *domain.Geography,
	pq *utils.PaginationQuery) (*domain.SpotList, error) {

	if err := utils.ValidateStruct(ctx, coordinate); err != nil {
		return nil, httpErrors.NewBadRequestError(errors.WithMessage(err, "spotService.GetSpotsInRangeOfCoordinate.ValidateCoordinate"))
	}

	allSpots, err := sps.spotRepo.GetSpotsInRange(ctx, coordinate, pq)

	if err != nil {
		return nil, err
	}
	//TODO: If 2 spots less than 50m, order by rate.

	return allSpots, nil

}
