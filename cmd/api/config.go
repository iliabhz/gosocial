package main

type appConfig struct {
	addr string
}

func newAppConfig(addr string) *appConfig {
	return &appConfig{
		addr: addr,
	}
}
