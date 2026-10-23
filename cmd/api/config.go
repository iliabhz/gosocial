package main

import "github.com/iliabhz/gosocial/internal/env"

type appConfig struct {
	addr string
}

func newAppConfig() *appConfig {
	addr := env.GetString("ADDR", ":8080")
	return &appConfig{
		addr: addr,
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
