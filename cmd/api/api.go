package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/iliabhz/gosocial/docs" // For swagger docs
	"github.com/iliabhz/gosocial/internal/store"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"go.uber.org/zap"
)

type application struct {
	mux    http.Handler
	store  store.Storage
	logger *zap.SugaredLogger
	env    string
}

func newApplication(cfg *appConfig, store store.Storage, logger *zap.SugaredLogger) *application {
	return &application{
		store:  store,
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

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", app.healthCheckHandler)

		docsURL := fmt.Sprintf("%s/swagger/doc.json", appCfg.addr)
		r.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL(docsURL)))

		//Public routes
		r.Route("/authentication", func(r chi.Router) {
			r.Post("/user", app.registerUserHandler(mailCfg))
		})

		r.Route("/posts", func(r chi.Router) {
			r.Post("/", app.createPostHandler)

			r.Route("/{postID}", func(r chi.Router) {
				r.Use(app.PostCtx)
				r.Get("/", app.getPostHandler)
				r.Patch("/", app.updatePostHandler)
				r.Delete("/", app.deletePostHandler)

				r.Route("/comments", func(r chi.Router) {
					r.Post("/", app.createCommentHandler)
				})
			})
		})

		r.Route("/users", func(r chi.Router) {
			r.Put("/activate", app.activateUserHandler)

			r.Route("/{userID}", func(r chi.Router) {
				r.Use(app.UserCtx)
				r.Get("/", app.getUserHandler)

				r.Put("/follow", app.followUserHandler)
				r.Put("/unfollow", app.unfollowUserHandler)
			})

			r.Group(func(r chi.Router) {
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

	app.logger.Infow("Server has started", "addr", cfg.addr, "env", app.env)

	return srv.ListenAndServe()
}
