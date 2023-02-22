package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/kkr2/spots/internal/spots/domain"
	"github.com/kkr2/spots/pkg/utils"
	"github.com/pkg/errors"
)

// SpotRepository is a repository interface
type SpotRepository interface {
	GetSpotsInRange(ctx context.Context, coordinate *domain.Geography, pq *utils.PaginationQuery) (*domain.SpotList, error)
}

// Spots Repository implementation 
type spotRepo struct {
	db *sqlx.DB
}

// NewSpotRepository initiates new repository of spots 
func NewSpotRepository(db *sqlx.DB) SpotRepository {
	return &spotRepo{db: db}
}

// GetSpotsInRange returns all geo spots in the given range of the coordinate 
func (sr *spotRepo) GetSpotsInRange(
	ctx context.Context,
	coordinate *domain.Geography,
	query *utils.PaginationQuery) (*domain.SpotList, error) {

	var totalCount int
	if err := sr.db.GetContext(ctx, &totalCount, findTotalSpotsInRange, coordinate.Latitude, coordinate.Longitude, query.Range); err != nil {
		return nil, errors.Wrap(err, "spotsRepo.GetTotalSpotsInRange.GetContext")
	}

	if totalCount == 0 {
		return &domain.SpotList{
			TotalCount: totalCount,
			TotalPages: utils.GetTotalPages(totalCount, query.GetSize()),
			Page:       query.GetPage(),
			Size:       query.GetSize(),
			HasMore:    utils.GetHasMore(query.GetPage(), totalCount, query.GetSize()),
			Spots:      make([]*domain.Spot, 0),
		}, nil
	}
		
	var spotList = make([]*domain.Spot, 0, query.GetSize())

	rows, err := sr.db.QueryxContext(ctx, findSpotsInRange, coordinate.Latitude, coordinate.Longitude, query.Range, query.GetOffset(), query.GetLimit())
	if err != nil {
		return nil, errors.Wrap(err, "spotsRepo.GetSpotsInRange.QueryxContext")
	}

	for rows.Next() {
		n := &domain.Spot{}
		if err = rows.StructScan(n); err != nil {
			return nil, errors.Wrap(err, "spotsRepo.GetSpotsInRange.StructScan")
		}
		spotList = append(spotList, n)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "spotsRepo.GetSpotsInRange.rows.Err")
	}

	return &domain.SpotList{
		TotalCount: totalCount,
		TotalPages: utils.GetTotalPages(totalCount, query.GetSize()),
		Page:       query.GetPage(),
		Size:       query.GetSize(),
		HasMore:    utils.GetHasMore(query.GetPage(), totalCount, query.GetSize()),
		Spots:      spotList,
	}, nil
}
