package main

import (
	"time"

	"github.com/iliabhz/gosocial/internal/env"
)

type appConfig struct {
	addr   string
	apiURL string
	env    string
}

func newAppConfig() *appConfig {
	addr := env.GetString("ADDR", ":8080")
	apiURL := env.GetString("EXTERNAL_URL", "localhost:8080")
	env := env.GetString("ENV", "dev")
	return &appConfig{
		addr:   addr,
		apiURL: apiURL,
		env:    env,
	}
}

type dbConfig struct {
	dsn string
}

func newDBConfig() *dbConfig {
	dsn := env.GetString("DB_DSN", "postgres://admin:adminpass@localhost/gosocial?sslmode=disable&pool_max_conns=30&pool_max_conn_idle_time=15m")
	return &dbConfig{
		dsn: dsn,
	}
}

type mailConfig struct {
	expiry time.Duration
}

func newMailConfig() *mailConfig {
	expiry, err := time.ParseDuration(env.GetString("MAIL_EXPIRY", "72h"))
	if err != nil {
		expiry = time.Hour * 24 * 3
	}

	return &mailConfig{
		expiry: expiry,
	}
}
