package config

import (
	"os"

	"golang.org/x/mod/modfile"
)

type AppEnv string

// AppConfig common app configs
type AppConfig struct {
	// App mode: local, develop, test, production
	AppEnv AppEnv `env:"APP_ENV" envDefault:"develop"`
	// uuid for use as namespace in uuid_v5
	NamespaceUUID string `env:"NAMESPACE_UUID" envDefault:"044754d5-e8ec-4821-95e5-74036610baa7"`
}

func (s *AppConfig) GetServiceName(goModPath string) (string, error) {
	path := "go.mod"
	if goModPath != "" {
		path = goModPath
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	modFile, err := modfile.Parse(path, data, nil)
	if err != nil {
		return "", err
	}

	return modFile.Module.Mod.Path, nil
}
