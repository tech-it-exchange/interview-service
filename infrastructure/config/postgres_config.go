package config

import (
	"fmt"
	"time"
)

type PostgresConfig struct {
	DatasourceUser            string        `env:"DATASOURCE_USER"`
	DatasourcePassword        string        `env:"DATASOURCE_PASSWORD"`
	DatasourceHost            string        `env:"DATASOURCE_HOST"`
	DatasourcePort            string        `env:"DATASOURCE_PORT"`
	DatasourceDatabase        string        `env:"DATASOURCE_DATABASE"`
	DatasourceMaxConns        int32         `env:"DATASOURCE_MAX_CONNS"          envDefault:"50"`
	DatasourceMaxLifetime     time.Duration `env:"DATASOURCE_MAX_LIFETIME"       envDefault:"3s"`
	DatasourceMinConns        int32         `env:"DATASOURCE_MIN_CONNS"          envDefault:"10"`
	DatasourceMaxConnIdleTime time.Duration `env:"DATASOURCE_MAX_CONN_IDLE_TIME" envDefault:"15s"`
	DatasourceConnectTimeout  time.Duration `env:"DATASOURCE_CONNECT_TIMEOUT"    envDefault:"5s"`
	DBSSL                     bool          `env:"DB_SSL"                        envDefault:"false"`
	DBSSLExtra                bool          `env:"DB_SSL_EXTRA"                  envDefault:"false"`
	datasourceUrl             string
}

func (s *PostgresConfig) GetDatasourceUrl() string {
	if s.datasourceUrl == "" {
		s.datasourceUrl = fmt.Sprintf(
			"postgresql://%s:%s@%s:%s/%s",
			s.DatasourceUser,
			s.DatasourcePassword,
			s.DatasourceHost,
			s.DatasourcePort,
			s.DatasourceDatabase,
		)
	}

	return s.datasourceUrl
}
