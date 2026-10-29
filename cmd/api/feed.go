package main

import (
	"net/http"

	"github.com/iliabhz/gosocial/internal/store"
)

func (app *application) getUserFeedHandler(w http.ResponseWriter, r *http.Request) {
	pf := store.PaginatedFeedQuery{
		Limit:  20,
		Offset: 0,
		Sort:   store.DESC_SORT,
	}

	pf, err := pf.Parse(r)
	if err != nil {
		writeJSONError(w, r, NewHTTPError("Invalid pagination query params", http.StatusBadRequest))
		return
	}

	if err := Validate.Struct(pf); err != nil {
		writeJSONError(w, r, err)
		return
	}

	ctx := r.Context()

	var userID int64 = 16 // Mock userID before auth
	feed, err := app.store.Posts().GetUserFeed(ctx, userID, pf)
	if err != nil {
		writeJSONError(w, r, err)
		return
	}

	writeDataJSON(w, http.StatusOK, feed)
}
