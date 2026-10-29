package store

import (
	"net/http"
	"strconv"
	"strings"
)

var (
	ASC_SORT  = "ASC"
	DESC_SORT = "DESC"
)

// Limit Lower Bound = 1 / Upper Bound = 20
// Offset Lower Bount = 0
// SortType: ASC / DESC

type PaginatedFeedQuery struct {
	Limit  int    `json:"limit" validate:"gte=1,lte=20"`
	Offset int    `json:"offset" validate:"gte=0"`
	Sort   string `json:"sort" validate:"oneof=ASC DESC"`
}

func (fq PaginatedFeedQuery) Parse(r *http.Request) (PaginatedFeedQuery, error) {
	queryStr := r.URL.Query()

	limit := queryStr.Get("limit")
	if limit != "" {
		l, err := strconv.Atoi(limit)
		if err != nil {
			return fq, nil
		}

		fq.Limit = l
	}

	offset := queryStr.Get("offset")
	if offset != "" {
		o, err := strconv.Atoi(offset)
		if err != nil {
			return fq, nil
		}

		fq.Offset = o
	}

	sort := queryStr.Get("sort")
	if sort != "" {
		fq.Sort = strings.ToUpper(sort)
	}

	return fq, nil
}
