package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/iliabhz/gosocial/internal/store"
)

const userContextKey = "user"
const postContextKey = "post"

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

		ctx := context.WithValue(r.Context(), postContextKey, post)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

type AccessCheck func(*http.Request) (bool, error)

func (app *application) RoleAtLeast(minRole string) AccessCheck {
	return func(r *http.Request) (bool, error) {
		user, ok := r.Context().Value(userContextKey).(*store.User)
		if !ok {
			return false, errors.New("failed to get user from context")
		}

		role, err := app.store.Roles().GetByName(r.Context(), minRole)
		if err != nil {
			return false, err
		}

		return user.Role.Level >= role.Level, nil
	}
}

func (app *application) PostOwner() AccessCheck {
	return func(r *http.Request) (bool, error) {
		ctx := r.Context()

		user, ok := ctx.Value(userContextKey).(*store.User)
		if !ok {
			return false, errors.New("failed to get user from context")
		}

		post, ok := ctx.Value(postContextKey).(*store.Post)
		if !ok {
			return false, errors.New("failed to get post from context")
		}

		return user.ID == post.UserID, nil
	}
}

func (app *application) RequireAny(checks ...AccessCheck) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			for _, check := range checks {
				allowed, err := check(r)

				if err != nil {
					writeJSONError(app.logger, w, r, err)
					return
				}

				if allowed {
					next.ServeHTTP(w, r)
					return
				}
			}

			writeJSONError(app.logger, w, r, ErrForbidden)
		})
	}
}

func (app *application) RequireAll(checks ...AccessCheck) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			for _, check := range checks {
				allowed, err := check(r)

				if err != nil {
					writeJSONError(app.logger, w, r, err)
					return
				}

				if !allowed {
					writeJSONError(app.logger, w, r, ErrForbidden)
					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}
