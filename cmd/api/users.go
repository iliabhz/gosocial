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
			writeJSONError(app.logger, w, r, NewHTTPError("Invalid user id", http.StatusBadRequest))
			return
		}
		user, err := app.store.Users().GetByID(r.Context(), userID)
		if err != nil {
			writeJSONError(app.logger, w, r, err)
			return
		}

		ctx := context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUser godoc
//
//	@Summary		Fetches a user profile
//	@Description	Fetches a user profile by ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"User ID"
//	@Success		200	{object}	store.User
//	@Failure		400	{object}	error
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/users/{id} [get]
func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, ok := ctx.Value("user").(*store.User)
	if !ok {
		writeJSONError(app.logger, w, r, NewHTTPError("Couldn't process the user", http.StatusUnprocessableEntity))
		return
	}

	writeDataJSON(w, http.StatusFound, user)
}

// FollowUser godoc
//
//	@Summary		Follows a user
//	@Description	Follows a user by ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			userID	path		int		true	"User ID"
//	@Success		204		{string}	string	"User followed"
//	@Failure		400		{object}	error	"User payload missing"
//	@Failure		404		{object}	error	"User not found"
//	@Failure		442		{object}	error	"Couldn't process the user"
//	@Security		ApiKeyAuth
//	@Router			/users/{userID}/follow [put]
func (app *application) followUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	followerUser, ok := ctx.Value("user").(*store.User)
	if !ok {
		writeJSONError(app.logger, w, r, NewHTTPError("Couldn't process the user", http.StatusUnprocessableEntity))
		return
	}

	var userID int64 = 16 // Mock userID before auth
	if err := app.store.Followers().Follow(ctx, followerUser.ID, userID); err != nil {
		writeJSONError(app.logger, w, r, err)
		return
	}

	writeJSON(w, http.StatusOK, "Followed successfully")
}

// UnfollowUser gdoc
//
//	@Summary		Unfollow a user
//	@Description	Unfollow a user by ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			userID	path		int		true	"User ID"
//	@Success		204		{string}	string	"User unfollowed"
//	@Failure		400		{object}	error	"User payload missing"
//	@Failure		404		{object}	error	"User not found"
//	@Security		ApiKeyAuth
//	@Router			/users/{userID}/unfollow [put]
func (app *application) unfollowUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	unfollowedUser, ok := ctx.Value("user").(*store.User)
	if !ok {
		writeJSONError(app.logger, w, r, NewHTTPError("Couldn't process the user", http.StatusUnprocessableEntity))
		return
	}

	var userID int64 = 4 // Mock userID before auth
	if err := app.store.Followers().Unfollow(ctx, unfollowedUser.ID, userID); err != nil {
		writeJSONError(app.logger, w, r, err)
		return
	}

	writeJSON(w, http.StatusOK, "User unfollowed")
}
