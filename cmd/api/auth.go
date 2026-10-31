package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"uuid"

	"github.com/iliabhz/gosocial/internal/mailer"
	"github.com/iliabhz/gosocial/internal/store"
)

type RegisterUserPayload struct {
	Username string `json:"username" validate:"required,max=100"`
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=3,max=72"`
}

// registerUserHandler godoc
//
//	@Summary		Registers a user
//	@Description	Registers a user
//	@Tags			authentication
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		RegisterUserPayload	true	"User credentials"
//	@Success		201		{object}	string				"User registered"
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Router			/authentication/user [post]
func (app *application) registerUserHandler(appCfg *appConfig, mailCfg *mailConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		payload := RegisterUserPayload{}
		if err := readJSON(w, r, &payload); err != nil {
			writeJSONError(app.logger, w, r, NewHTTPError("Invalid json", http.StatusBadRequest))
			return
		}

		if err := Validate.Struct(&payload); err != nil {
			writeJSONError(app.logger, w, r, err)
			return
		}
		ctx := r.Context()

		user := &store.User{
			Username: payload.Username,
			Email:    payload.Email,
		}

		if err := user.Password.Set(payload.Password); err != nil {
			writeJSONError(app.logger, w, r, err)
			return
		}

		plainToken := uuid.NewV4().String()

		hash := sha256.Sum256([]byte(plainToken))
		hashToken := hex.EncodeToString(hash[:])

		// add config for expiry
		if err := app.store.Users().CreateAndInvite(ctx, mailCfg.expiry, hashToken, user); err != nil {
			writeJSONError(app.logger, w, r, err)
			return
		}

		// DEV ONLY
		if app.env == "dev" {
			app.logger.Infow("Registered user", "id", user.ID, "username", user.Username, "email", user.Email, "token", plainToken)
		}

		activationURL := fmt.Sprintf("%s/confirm/%s", appCfg.frontendURL, plainToken)

		isSandbox := app.env == "dev"
		vars := struct {
			Username      string
			ActivationURL string
		}{
			Username:      user.Username,
			ActivationURL: activationURL,
		}

		err := app.mailer.Send(mailer.UserInvitationTemplate, user.Username, user.Email, vars, isSandbox)
		if err != nil {
			app.logger.Errorw("error sending user invitation email", "error", err.Error())

			// Rollback user creation
			if err := app.store.Users().Delete(ctx, user.ID); err != nil {
				app.logger.Errorw("error deleting user after failing sending an email", "error", err.Error())
				writeJSONError(app.logger, w, r, err)
				return
			}

			writeJSONError(app.logger, w, r, err)
			return
		}

		writeJSON(w, http.StatusCreated, "User registered")
	}
}
