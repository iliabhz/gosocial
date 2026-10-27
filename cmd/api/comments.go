package main

import (
	"net/http"

	"github.com/iliabhz/gosocial/internal/store"
)

type CreateCommentPayload struct {
	Content string `json:"content" validate:"required,max=1000"`
}

func (app *application) createCommentHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreateCommentPayload
	if err := readJSON(w, r, &payload); err != nil {
		writeJSONError(w, r, NewHTTPError("Invalid json", http.StatusBadRequest))
		return
	}

	if err := Validate.Struct(payload); err != nil {
		writeJSONError(w, r, err)
		return
	}

	ctx := r.Context()

	post, ok := ctx.Value("post").(*store.Post)
	if !ok {
		writeJSONError(w, r, NewHTTPError("Couldn't process the post", http.StatusUnprocessableEntity))
		return
	}

	user := store.User{ // Mock the user before auth for now
		ID:       1,
		Username: "bob",
	}

	comment := store.Comment{
		PostID:  post.ID,
		UserID:  user.ID,
		Content: payload.Content,
		User:    user,
	}

	if err := app.store.Comments().Create(ctx, &comment); err != nil {
		writeJSONError(w, r, err)
		return
	}

	if err := writeDataJSON(w, http.StatusCreated, &comment); err != nil {
		writeJSONError(w, r, err)
		return
	}
}
