package main

import (
	"net/http"

	"github.com/iliabhz/gosocial/internal/store"
)

type CreatePostPayload struct {
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
}

func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreatePostPayload
	if err := readJSON(w, r, &payload); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	ctx := r.Context()

	post := store.Post{
		Title:   payload.Title,
		Content: payload.Content,
		Tags:    payload.Tags,
		UserID:  1, // Mock userid before auth
	}

	if err := app.store.Posts().Create(ctx, &post); err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := writeJSON(w, http.StatusCreated, &post); err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
}
