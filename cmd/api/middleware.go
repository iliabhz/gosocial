package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/iliabhz/gosocial/internal/store"
)

func (app *application) JWTAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			writeJSONError(app.logger, w, r, ErrUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			writeJSONError(app.logger, w, r, ErrUnauthorized)
			return
		}

		token := parts[1]

		jwtToken, err := JWTAuth.ValidateToken(token)
		if err != nil {
			writeJSONError(app.logger, w, r, ErrUnauthorized)
			return
		}

		claims, ok := jwtToken.Claims.(jwt.MapClaims)
		if !ok {
			writeJSONError(app.logger, w, r, ErrUnauthorized)
			return
		}

		userID, err := strconv.ParseInt(fmt.Sprintf("%.f", claims["sub"]), 10, 64)
		if err != nil {
			writeJSONError(app.logger, w, r, ErrUnauthorized)
			return
		}

		ctx := r.Context()
		user, err := app.store.Users().GetByID(ctx, userID)
		if err != nil {
			writeJSONError(app.logger, w, r, ErrUnauthorized)
			return
		}

		ctx = context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (app *application) PostCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		postIDStr := chi.URLParam(r, "postID")
		postID, err := strconv.ParseInt(postIDStr, 10, 64)
		if err != nil {
			writeJSONError(app.logger, w, r, NewHTTPError("Invalid post id", http.StatusBadRequest))
			return
		}
		post, err := app.store.Posts().GetByID(r.Context(), postID)
		if err != nil {
			writeJSONError(app.logger, w, r, err)
			return
		}

		ctx := context.WithValue(r.Context(), "post", post)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (app *application) RolePrecedenceMiddleware(roleName string, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		role, err := app.store.Roles().GetByName(ctx, roleName)
		if err != nil {
			writeJSONError(app.logger, w, r, err)
			return
		}
		user, _ := ctx.Value("user").(*store.User)

		if user.Role.Level >= role.Level {
			next.ServeHTTP(w, r)
			return
		}

		writeJSONError(app.logger, w, r, ErrForbidden)
	})
}

func (app *application) PostOwnershipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		user, _ := ctx.Value("user").(*store.User)
		post, _ := ctx.Value("post").(*store.Post)

		if user.ID == post.UserID {
			next.ServeHTTP(w, r)
			return
		}

	})
}
