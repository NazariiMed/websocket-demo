package models

import "github.com/caarlos0/env"

type Config struct {
	DatabaseUrl string `env:"DATABASE_URL"    envDefault:"postgres://username:password@localhost:5432/websocket_demo"`
	RestPort    int    `env:"REST_PORT"       envDefault:"8080"`
	WsUrl       string `env:"WS_URL"          envDefault:"wss://stream.binance.com:9443/ws"`
	HttpBaseUrl string `env:"HTTP_BASE_URL"   envDefault:"https://vapi.binance.com/"`
}

func ParseConfig() (Config, error) {
	var out Config
	err := env.Parse(&out)
	return out, err
}
