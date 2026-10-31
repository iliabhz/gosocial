package main

import (
	"time"

	"github.com/iliabhz/gosocial/internal/env"
)

type appConfig struct {
	addr        string
	apiURL      string
	frontendURL string
	env         string
}

func newAppConfig() *appConfig {
	addr := env.GetString("ADDR", ":8080")
	apiURL := env.GetString("EXTERNAL_URL", "localhost:8080")
	frontendURL := env.GetString("FRONTEND_URL", "localhost:4000")
	env := env.GetString("ENV", "dev")
	return &appConfig{
		addr:        addr,
		apiURL:      apiURL,
		frontendURL: frontendURL,
		env:         env,
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
	mailGun   *mailGunConfig
	domain    string
	fromEmail string
	expiry    time.Duration
}

type mailGunConfig struct {
	apiKey string
}

func newMailConfig() *mailConfig {
	expiry, err := time.ParseDuration(env.GetString("MAIL_EXPIRY", "72h"))
	if err != nil {
		expiry = time.Hour * 24 * 3
	}

	mgCfg := mailGunConfig{
		apiKey: env.GetString("MAILGUN_API_KEY", ""),
	}

	return &mailConfig{
		mailGun:   &mgCfg,
		domain:    env.GetString("MAIL_DOMAIN", ""),
		fromEmail: env.GetString("MAIL_FROM_EMAIL", ""),
		expiry:    expiry,
	}
}
