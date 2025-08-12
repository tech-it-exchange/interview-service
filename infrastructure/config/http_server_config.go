package config

type HttpServerConfig struct {
	HttpPort string `env:"SERVICE_PORT" envDefault:"8080"`
}
