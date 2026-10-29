package main

import (
	"net/http"

	"github.com/iliabhz/gosocial/internal/store"
)

// getUserFeedHandler godoc
//
//	@Summary		Fetches the user feed
//	@Description	Fetches the user feed
//	@Tags			feed
//	@Accept			json
//	@Produce		json
//	@Param			since	query		string	false	"Since"
//	@Param			until	query		string	false	"Until"
//	@Param			limit	query		int		false	"Limit"
//	@Param			offset	query		int		false	"Offset"
//	@Param			sort	query		string	false	"Sort"
//	@Param			tags	query		string	false	"Tags"
//	@Param			search	query		string	false	"Search"
//	@Success		200		{object}	[]store.FeedPost
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/users/feed [get]
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

	ff := store.FilterdFeedQuery{}

	ff, err = ff.Parse(r)
	if err != nil {
		writeJSONError(w, r, NewHTTPError("Invalid filter query params", http.StatusBadRequest))
		return
	}

	if err := Validate.Struct(ff); err != nil {
		writeJSONError(w, r, err)
		return
	}

	if !ff.Since.IsZero() && !ff.Until.IsZero() && ff.Since.After(ff.Until) {
		writeJSONError(w, r, NewHTTPError("since date is after until date", http.StatusBadRequest))
		return
	}

	ctx := r.Context()

	var userID int64 = 16 // Mock userID before auth
	feed, err := app.store.Posts().GetUserFeed(ctx, userID, pf, ff)
	if err != nil {
		writeJSONError(w, r, err)
		return
	}

	writeDataJSON(w, http.StatusOK, feed)
}
