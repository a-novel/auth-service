package config

import (
	"context"
	_ "embed"
	"github.com/a-novel/auth-service/migrations"
	"github.com/a-novel/go-framework/postgresql/bunframework"
	"github.com/a-novel/go-framework/postgresql/bunframework/pgconfig"
	"github.com/rs/zerolog"
	"github.com/uptrace/bun"
	"io/fs"
	"log"
	"time"
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

func GetPostgres(logger zerolog.Logger) (*bun.DB, func()) {
	postgres, sql, err := bunframework.NewClient(context.Background(), bunframework.Config{
		Driver: pgconfig.Driver{
			DSN:         Postgres.DSN,
			AppName:     App.Name,
			DialTimeout: 120 * time.Second,
		},
		DiscardUnknownColumns: true,
		Migrations: &bunframework.MigrateConfig{
			Files: []fs.FS{migrations.Migrations},
		},
	})

	if err != nil {
		logger.Fatal().Err(err).Msg("error connecting to postgres")
	}

	return postgres, func() {
		_ = postgres.Close()
		_ = sql.Close()
	}
}
