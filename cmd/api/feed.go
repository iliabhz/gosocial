package main

import "net/http"

func (app *application) getUserFeedHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var userID int64 = 16 // Mock userID before auth
	feed, err := app.store.Posts().GetUserFeed(ctx, userID)
	if err != nil {
		writeJSONError(w, r, err)
		return
	}

	writeDataJSON(w, http.StatusOK, feed)
}
