package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/kkr2/spots/internal/spots/domain"
	"github.com/kkr2/spots/pkg/utils"
	"github.com/pkg/errors"
)

type SpotRepository interface {
	GetSpotsInRange(ctx context.Context, coordinate *domain.Geography, pq *utils.PaginationQuery) (*domain.SpotList, error)
}

// News Repository
type spotRepo struct {
	db *sqlx.DB
}

// News repository constructor
func NewNewsRepository(db *sqlx.DB) SpotRepository {
	return &spotRepo{db: db}
}

func (sr *spotRepo) GetSpotsInRange(
	ctx context.Context,
	coordinate *domain.Geography,
	query *utils.PaginationQuery) (*domain.SpotList, error) {

	var totalCount int
	if err := sr.db.GetContext(ctx, &totalCount, findTotalSpotsInRange,coordinate.Latitude,coordinate.Longitude, query.Range); err != nil {
		return nil, errors.Wrap(err, "spotRepo.GetTotalSpotsInRange.GetContext")
	}
	if totalCount == 0 {
		return &domain.SpotList{
			TotalCount: totalCount,
			TotalPages: utils.GetTotalPages(totalCount, query.GetSize()),
			Page:       query.GetPage(),
			Size:       query.GetSize(),
			HasMore:    utils.GetHasMore(query.GetPage(), totalCount, query.GetSize()),
			Spots:       make([]*domain.Spot, 0),
		}, nil
	}

	var spotList = make([]*domain.Spot, 0, query.GetSize())
	rows, err := sr.db.QueryxContext(ctx, findSpotsInRange, coordinate.Latitude,coordinate.Longitude, query.Range, query.GetOffset(), query.GetLimit())
	if err != nil {
		return nil, errors.Wrap(err, "spotRepo.GetSpotsInRange.QueryxContext")
	}
	defer rows.Close()

	for rows.Next() {
		n := &domain.Spot{}
		if err = rows.StructScan(n); err != nil {
			return nil, errors.Wrap(err, "newsRepo.GetSpotsInRange.StructScan")
		}
		spotList = append(spotList, n)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "newsRepo.GetSpotsInRange.rows.Err")
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
