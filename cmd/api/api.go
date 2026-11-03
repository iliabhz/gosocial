package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/iliabhz/gosocial/docs" // For swagger docs
	"github.com/iliabhz/gosocial/internal/env"
	"github.com/iliabhz/gosocial/internal/mailer"
	"github.com/iliabhz/gosocial/internal/store"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"go.uber.org/zap"
)

type application struct {
	mux    http.Handler
	store  store.Storage
	mailer mailer.Client
	logger *zap.SugaredLogger
	env    string
}

func newApplication(cfg *appConfig, store store.Storage, mailer mailer.Client, logger *zap.SugaredLogger) *application {
	return &application{
		store:  store,
		mailer: mailer,
		logger: logger,
		env:    cfg.env,
	}
}

func (app *application) mount(appCfg *appConfig, mailCfg *mailConfig) {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.ClientIPFromRemoteAddr)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	allowedOrigins := env.GetStringArray("CORS_ALLOWED_ORIGINS", []string{env.GetString("FRONTEND_URL", "http://localhost:4000")})
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/v1", func(r chi.Router) {
		r.With(app.JWTAuthMiddleware).Get("/health", app.healthCheckHandler)

		docsURL := fmt.Sprintf("%s/swagger/doc.json", appCfg.addr)
		r.With(app.JWTAuthMiddleware).Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL(docsURL)))

		//Public routes
		r.Route("/authentication", func(r chi.Router) {
			r.Post("/user", app.registerUserHandler(appCfg, mailCfg))
			r.Post("/token", app.createTokenHandler)
		})

		r.Route("/posts", func(r chi.Router) {
			r.Use(app.JWTAuthMiddleware)

			r.Post("/", app.createPostHandler)

			r.Route("/{postID}", func(r chi.Router) {
				r.Use(app.PostCtx)
				r.Get("/", app.getPostHandler)

				r.With(app.RequireAny(
					app.PostOwner(),
					app.RoleAtLeast("moderator"),
				)).Patch("/", app.updatePostHandler)

				r.With(app.RequireAny(
					app.PostOwner(), app.RoleAtLeast("admin"),
				)).Delete("/", app.deletePostHandler)

				r.Route("/comments", func(r chi.Router) {
					r.Post("/", app.createCommentHandler)
				})
			})
		})

		r.Route("/users", func(r chi.Router) {
			r.Put("/activate", app.activateUserHandler)

			r.Route("/{userID}", func(r chi.Router) {
				r.Use(app.JWTAuthMiddleware)

				r.Get("/", app.getUserHandler)

				r.Put("/follow", app.followUserHandler)
				r.Put("/unfollow", app.unfollowUserHandler)
			})

			r.Group(func(r chi.Router) {
				r.Use(app.JWTAuthMiddleware)
				r.Get("/feed", app.getUserFeedHandler)
			})
		})
	})

	app.mux = r
}

func (app *application) run(cfg *appConfig) error {
	// Docs
	docs.SwaggerInfo.Version = version
	docs.SwaggerInfo.Host = cfg.apiURL

	srv := &http.Server{
		Addr:         cfg.addr,
		Handler:      app.mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	shutdown := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)

		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		app.logger.Infow("signal caught", "signal", s.String())

		shutdown <- srv.Shutdown(ctx)
	}()

	app.logger.Infow("Server has started", "addr", cfg.addr, "env", app.env, "PID", os.Getpid())

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdown
	if err != nil {
		return err
	}

	app.logger.Infow("server has stopped", "addr", cfg.addr, "env", app.env)

	return nil
}
