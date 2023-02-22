package utils

import (
	"fmt"
	"math"
	"strconv"

	"github.com/labstack/echo/v4"
)

const (
	defaultSize  = 10
	defaultRange = 50
)

// PaginationQuery is a structure that holds information that repository uses for pagination
type PaginationQuery struct {
	Size  int `json:"size,omitempty"`
	Page  int `json:"page,omitempty"`
	Range int `json:"range,omitempty"`
}

// etSize sets query size to pagination object
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

// SetPage sets the page nr to pagination object
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

// SetRange sets desired range to pagination object
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

// GetOffset gets offset from pagination object
func (q *PaginationQuery) GetOffset() int {
	if q.Page == 0 {
		return 0
	}
	return (q.Page - 1) * q.Size
}

// GetLimit gets limit from pagination object
func (q *PaginationQuery) GetLimit() int {
	return q.Size
}

// GetRange gets range from pagination object
func (q *PaginationQuery) GetRange() int {
	return q.Range
}

// GetPage gets page nr from pagination object
func (q *PaginationQuery) GetPage() int {
	return q.Page
}

// GetSize returns size from pagination query
func (q *PaginationQuery) GetSize() int {
	return q.Size
}

// GetQueryString returns a stringifyed verson with pagination query params
func (q *PaginationQuery) GetQueryString() string {
	return fmt.Sprintf("page=%v&size=%v&range=%v", q.GetPage(), q.GetSize(), q.GetRange())
}

// GetPaginationFromCtx returns a pagination object given an echo context
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

// GetTotalPages calculates total pages given totalcount and pagesize
func GetTotalPages(totalCount int, pageSize int) int {
	d := float64(totalCount) / float64(pageSize)
	return int(math.Ceil(d))
}

// GetHasMore returns if db has more data 
func GetHasMore(currentPage int, totalCount int, pageSize int) bool {
	return currentPage < totalCount/pageSize
}
