package main

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/iliabhz/gosocial/internal/store"
)

type CreatePostPayload struct {
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
}

func (app *application) PostCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		postIDStr := chi.URLParam(r, "postID")
		postID, err := strconv.ParseInt(postIDStr, 10, 64)
		if err != nil {
			writeJSONError(w, http.StatusBadRequest, "Invalid post id")
			return
		}
		post, err := app.store.Posts().GetByID(r.Context(), postID)
		if err != nil {
			switch {
			case errors.Is(err, store.ErrPostNotFound):
				writeJSONError(w, http.StatusNotFound, err.Error())
			default:
				writeJSONError(w, http.StatusInternalServerError, err.Error())
			}

			return
		}
		ctx := context.WithValue(r.Context(), "post", post)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (app *application) getPostHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	post, ok := ctx.Value("post").(*store.Post)
	if !ok {
		writeJSONError(w, http.StatusUnprocessableEntity, "Couldn't process the post")
	}

	writeJSON(w, http.StatusFound, post)
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
