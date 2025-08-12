package config

type ManagementConfig struct {
	ManagementPort string `env:"MANAGEMENT_PORT" envDefault:"9090"`
}
