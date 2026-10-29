package main

import "net/http"

// healthcheckHandler godoc
//
//	@Summary		Healthcheck
//	@Description	Healthcheck endpoint
//	@Tags			ops
//	@Produce		json
//	@Success		200	{object}	string	"ok"
//	@Router			/health [get]
func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status":  "ok",
		"env":     app.env,
		"version": version,
	}

	if err := writeJSON(w, http.StatusOK, data); err != nil {
		writeJSONError(app.logger, w, r, err)
	}
}
