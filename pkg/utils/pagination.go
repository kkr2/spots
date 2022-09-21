package utils

import (
	"fmt"
	"math"
	"strconv"

	"github.com/labstack/echo/v4"
)

const (
	defaultSize = 10
	defaultRange = 50
)

// Pagination query params
type PaginationQuery struct {
	Size  int `json:"size,omitempty"`
	Page  int `json:"page,omitempty"`
	Range int `json:"range,omitempty"`
}

// Set page size
func (q *PaginationQuery) SetSize(sizeQuery string) error {
	if sizeQuery == "" {
		q.Size = defaultSize
		return nil
	}
	n, err := strconv.Atoi(sizeQuery)
	if err != nil {
		return err
	}
	q.Size = n

	return nil
}

// Set page number
func (q *PaginationQuery) SetPage(pageQuery string) error {
	if pageQuery == "" {
		q.Size = 0
		return nil
	}
	n, err := strconv.Atoi(pageQuery)
	if err != nil {
		return err
	}
	q.Page = n

	return nil
}

// Set by range
func (q *PaginationQuery) SetRange(desiredRange string) error {
	if desiredRange == "" {
		q.Range = defaultRange
		return nil
	}
	n, err := strconv.Atoi(desiredRange)
	if err != nil {
		return err
	}
	q.Range = n

	return nil
}

// Get offset
func (q *PaginationQuery) GetOffset() int {
	if q.Page == 0 {
		return 0
	}
	return (q.Page - 1) * q.Size
}

// Get limit
func (q *PaginationQuery) GetLimit() int {
	return q.Size
}

// Get OrderBy
func (q *PaginationQuery) GetRange() int {
	return q.Range
}

// Get OrderBy
func (q *PaginationQuery) GetPage() int {
	return q.Page
}

// Get OrderBy
func (q *PaginationQuery) GetSize() int {
	return q.Size
}

func (q *PaginationQuery) GetQueryString() string {
	return fmt.Sprintf("page=%v&size=%v&range=%v", q.GetPage(), q.GetSize(), q.GetRange())
}

// Get pagination query struct from
func GetPaginationFromCtx(c echo.Context) (*PaginationQuery, error) {
	q := &PaginationQuery{}
	if err := q.SetPage(c.QueryParam("page")); err != nil {
		return nil, err
	}
	if err := q.SetSize(c.QueryParam("size")); err != nil {
		return nil, err
	}
	if err := q.SetRange(c.QueryParam("range")); err != nil {
		return nil, err
	}

	return q, nil
}

// Get total pages int
func GetTotalPages(totalCount int, pageSize int) int {
	d := float64(totalCount) / float64(pageSize)
	return int(math.Ceil(d))
}

// Get has more
func GetHasMore(currentPage int, totalCount int, pageSize int) bool {
	return currentPage < totalCount/pageSize
}
