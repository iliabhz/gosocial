package main

import "github.com/iliabhz/gosocial/internal/env"

type appConfig struct {
	addr string
}

func newAppConfig() *appConfig {
	return &appConfig{}
}

func initAppConfig(cfg *appConfig) {
	cfg.addr = env.GetString("ADDR", ":8080")
}

type dbConfig struct {
	dsn string
}

func newDBConfig() *dbConfig {
	return &dbConfig{}
}

func initDBConfig(cfg *dbConfig) {
	cfg.dsn = env.GetString("DB_DSN", "postgres://admin:adminpass@localhost/gosocial?sslmode=disable&pool_max_conns=30&pool_max_conn_idle_time=15m")
}
