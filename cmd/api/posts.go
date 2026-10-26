package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/iliabhz/gosocial/internal/store"
)

type CreatePostPayload struct {
	Title   string   `json:"title" validate:"required,max=100"`
	Content string   `json:"content" validate:"required,max=1000"`
	Tags    []string `json:"tags"`
}

type UpdatePostPayload struct {
	Title   *string  `json:"title" validate:"omitempty,max=100"`
	Content *string  `json:"content" validate:"omitempty,max=1000"`
	Tags    []string `json:"tags"`
}

func (app *application) PostCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		postIDStr := chi.URLParam(r, "postID")
		postID, err := strconv.ParseInt(postIDStr, 10, 64)
		if err != nil {
			writeJSONError(w, r, NewHTTPError("Invalid post id", http.StatusBadRequest))
			return
		}
		post, err := app.store.Posts().GetByID(r.Context(), postID)
		if err != nil {
			writeJSONError(w, r, err)
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
		writeJSONError(w, r, NewHTTPError("Couldn't process the post", http.StatusUnprocessableEntity))
		return
	}

	comments, err := app.store.Comments().GetByPostID(ctx, post.ID)
	if err != nil {
		writeJSONError(w, r, err)
		return
	}

	post.Comments = comments

	writeJSON(w, http.StatusFound, post)
}

func (app *application) updatePostHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	post, ok := ctx.Value("post").(*store.Post)
	if !ok {
		writeJSONError(w, r, NewHTTPError("Couldn't process the post", http.StatusUnprocessableEntity))
		return
	}

	var payload UpdatePostPayload
	if err := readJSON(w, r, &payload); err != nil {
		writeJSONError(w, r, NewHTTPError("Invalid json", http.StatusBadRequest))
		return
	}

	if err := Validate.Struct(payload); err != nil {
		writeJSONError(w, r, err)
		return
	}

	if payload.Title == nil && payload.Content == nil && len(payload.Tags) == 0 {
		writeJSONError(w, r, NewHTTPError("Invalid update payload", http.StatusBadRequest))
		return
	}

	if err := app.store.Posts().UpdateByID(ctx, payload.Title, payload.Content, payload.Tags, post.ID); err != nil {
		writeJSONError(w, r, err)
		return
	}

	writeJSON(w, http.StatusOK, "Updated the post successfully")
}

func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreatePostPayload
	if err := readJSON(w, r, &payload); err != nil {
		writeJSONError(w, r, NewHTTPError("Invalid json", http.StatusBadRequest))
		return
	}

	if err := Validate.Struct(payload); err != nil {
		writeJSONError(w, r, err)
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
		writeJSONError(w, r, err)
		return
	}

	if err := writeJSON(w, http.StatusCreated, &post); err != nil {
		writeJSONError(w, r, err)
		return
	}
}

func (app *application) deletePostHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	post, ok := ctx.Value("post").(*store.Post)
	if !ok {
		writeJSONError(w, r, NewHTTPError("Couldn't process the post", http.StatusUnprocessableEntity))
	}

	err := app.store.Posts().DeleteByID(ctx, post.ID)
	if err != nil {
		writeJSONError(w, r, err)
	}

	writeJSON(w, http.StatusOK, fmt.Sprintf("post with id: %v deleted successfully", post.ID))
}
