package main

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/iliabhz/gosocial/internal/store"
)

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
	followerUserIDStr := chi.URLParam(r, "userID")
	followerUserID, err := strconv.ParseInt(followerUserIDStr, 10, 64)
	if err != nil {
		writeJSONError(app.logger, w, r, NewHTTPError("Invalid user id", http.StatusBadRequest))
		return
	}

	ctx := r.Context()
	user, _ := ctx.Value("user").(*store.User)

	if err := app.store.Followers().Follow(ctx, followerUserID, user.ID); err != nil {
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
	unfollowedUserIDStr := chi.URLParam(r, "userID")
	unfollowedUserID, err := strconv.ParseInt(unfollowedUserIDStr, 10, 64)
	if err != nil {
		writeJSONError(app.logger, w, r, NewHTTPError("Invalid user id", http.StatusBadRequest))
		return
	}

	ctx := r.Context()
	user, _ := ctx.Value("user").(*store.User)
	if err := app.store.Followers().Unfollow(ctx, unfollowedUserID, user.ID); err != nil {
		writeJSONError(app.logger, w, r, err)
		return
	}

	writeJSON(w, http.StatusOK, "User unfollowed")
}

type ActivateUserPayload struct {
	Token string `json:"token" validate:"required,max=255"`
}

// ActivateUser godoc
//
//	@Summary		Activates/Register a user
//	@Description	Activates/Register a user by invitation token
//	@Tags			users
//	@Produce		json
//	@Param			token	body		ActivateUserPayload	true	"Invitation token"
//	@Success		204		{string}	string				"User activated"
//	@Failure		404		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/users/activate [put]
func (app *application) activateUserHandler(w http.ResponseWriter, r *http.Request) {
	payload := ActivateUserPayload{}
	if err := readJSON(w, r, &payload); err != nil {
		writeJSONError(app.logger, w, r, NewHTTPError("Invalid json", http.StatusBadRequest))
		return
	}

	if err := Validate.Struct(payload); err != nil {
		writeJSONError(app.logger, w, r, err)
		return
	}

	ctx := r.Context()

	if err := app.store.Users().Activate(ctx, payload.Token); err != nil {
		writeJSONError(app.logger, w, r, err)
		return
	}

	writeJSON(w, http.StatusOK, "User activated")
}
