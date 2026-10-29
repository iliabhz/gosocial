package store

import (
	"net/http"
	"strings"
	"time"
)

type FilterdFeedQuery struct {
	Search string    `json:"search" validate:"max=100"`
	Tags   []string  `json:"tags" validate:"max=5,dive,required,max=50"`
	Since  time.Time `json:"since"`
	Until  time.Time `json:"until"`
}

func (fq FilterdFeedQuery) Parse(r *http.Request) (FilterdFeedQuery, error) {
	queryStr := r.URL.Query()

	search := queryStr.Get("search")
	if search != "" {
		fq.Search = search
	}

	tags := queryStr.Get("tags")
	if tags != "" {
		fq.Tags = strings.Split(tags, ",")
	}

	since := queryStr.Get("since")
	if since != "" {
		sinceTime, err := time.Parse(time.RFC3339, since)
		if err != nil {
			return fq, nil
		}
		fq.Since = sinceTime
	}

	until := queryStr.Get("until")
	if until != "" {
		untilTime, err := time.Parse(time.RFC3339, until)
		if err != nil {
			return fq, nil
		}
		fq.Until = untilTime
	}

	return fq, nil
}
