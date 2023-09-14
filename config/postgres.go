package config

import (
	_ "embed"
	"log"
)

//go:embed postgres.yml
var postgresFile []byte

type PostgresConfig struct {
	DSN string `yaml:"dsn"`
}

var Postgres *PostgresConfig

func init() {
	cfg := new(PostgresConfig)

	if err := loadEnv(EnvLoader{DefaultENV: postgresFile}, cfg); err != nil {
		log.Fatalf("error loading postgres configuration: %v\n", err)
	}

	Postgres = cfg
}
