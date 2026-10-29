package main

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/iliabhz/gosocial/internal/store"
)

func (app *application) UserCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userIDStr := chi.URLParam(r, "userID")
		userID, err := strconv.ParseInt(userIDStr, 10, 64)
		if err != nil {
			writeJSONError(w, r, NewHTTPError("Invalid user id", http.StatusBadRequest))
			return
		}
		user, err := app.store.Users().GetByID(r.Context(), userID)
		if err != nil {
			writeJSONError(w, r, err)
			return
		}

		ctx := context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, ok := ctx.Value("user").(*store.User)
	if !ok {
		writeJSONError(w, r, NewHTTPError("Couldn't process the user", http.StatusUnprocessableEntity))
		return
	}

	writeDataJSON(w, http.StatusFound, user)
}

func (app *application) followUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	followerUser, ok := ctx.Value("user").(*store.User)
	if !ok {
		writeJSONError(w, r, NewHTTPError("Couldn't process the user", http.StatusUnprocessableEntity))
		return
	}

	var userID int64 = 16 // Mock userID before auth
	if err := app.store.Followers().Follow(ctx, followerUser.ID, userID); err != nil {
		writeJSONError(w, r, err)
		return
	}

	writeJSON(w, http.StatusOK, "Followed successfully")
}

func (app *application) unfollowUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	unfollowedUser, ok := ctx.Value("user").(*store.User)
	if !ok {
		writeJSONError(w, r, NewHTTPError("Couldn't process the user", http.StatusUnprocessableEntity))
		return
	}

	var userID int64 = 4 // Mock userID before auth
	if err := app.store.Followers().Unfollow(ctx, unfollowedUser.ID, userID); err != nil {
		writeJSONError(w, r, err)
		return
	}

	writeJSON(w, http.StatusOK, "Unfollowed successfully")
}
