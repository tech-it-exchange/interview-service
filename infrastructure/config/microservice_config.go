package config

import (
	"fmt"
	"github.com/caarlos0/env/v11"
	"runtime"
	"time"

	"github.com/joho/godotenv"
)

type MicroserviceConfig struct {
	ProcsNum                  int           `env:"PROCS_NUM" envDefault:"2"`
	RefreshInstrumentsTimeout time.Duration `env:"REFRESH_INSTRUMENTS_TIMEOUT" envDefault:"10s"`
	App                       AppConfig     `envPrefix:""`
	ChannelConfig             ChannelConfig `envPrefix:""`
	Management                ManagementConfig
	Swagger                   SwaggerConfig
	HttpServer                HttpServerConfig
}

var configuration *MicroserviceConfig

func init() {
	cfg := &MicroserviceConfig{}

	// Load environment variables from .env file
	if err := godotenv.Load(".env"); err != nil {
		fmt.Println("Error loading .env file, proceeding with environment variables only")
	}

	// Initialize the configuration
	if err := env.Parse(cfg); err != nil {
		fmt.Println("Error parsing environment variables:", err)
		return
	}

	configuration = cfg

	numProcs := runtime.GOMAXPROCS(0)
	if numProcs > configuration.ProcsNum {
		runtime.GOMAXPROCS(configuration.ProcsNum)
	}
}

// GetConfig returns the application configuration.
func GetConfig() *MicroserviceConfig {
	return configuration
}
