package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/iliabhz/gosocial/internal/store"
)

type application struct {
	mux   http.Handler
	store store.Storage
}

func newApplication(store store.Storage) *application {
	return &application{
		store: store,
	}
}

func (app *application) mount() {
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
	})

	app.mux = r
}

func (app *application) run(cfg *appConfig) error {
	srv := &http.Server{
		Addr:         cfg.addr,
		Handler:      app.mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	log.Printf("Server has started on %s", cfg.addr)

	return srv.ListenAndServe()
}
