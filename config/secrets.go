package config

import (
	"cloud.google.com/go/storage"
	"context"
	_ "embed"
	"github.com/a-novel/auth-service/pkg/dao"
	"github.com/rs/zerolog"
	"log"
	"os"
	"path"
	"time"
)

//go:embed secrets.yml
var secretsFile []byte

type SecretsConfig struct {
	Prefix         string        `yaml:"prefix"`
	Backups        int           `yaml:"backups"`
	UpdateInterval time.Duration `yaml:"updateInterval"`
}

var Secrets *SecretsConfig

func init() {
	cfg := new(SecretsConfig)

	if err := loadEnv(EnvLoader{DefaultENV: secretsFile}, cfg); err != nil {
		log.Fatalf("error loading secrets configuration: %v\n", err)
	}

	Secrets = cfg
}

func GetSecretsRepository(logger zerolog.Logger) (dao.SecretKeysRepository, zerolog.Logger) {
	if ENV == ProdENV {
		client, err := storage.NewClient(context.Background())
		if err != nil {
			logger.Fatal().Err(err).Msg("error initializing GCP client")
		}

		logger = logger.With().
			Dict(
				"secrets_manager",
				zerolog.Dict().
					Int("backups", Secrets.Backups).
					Dur("update_interval", Secrets.UpdateInterval).
					Str("type", "GCP Datastore"),
			).
			Logger()

		return dao.NewGoogleDatastoreSecretKeysRepository(client.Bucket(Deploy.Buckets.SecretKeys)), logger
	}

	wd, err := os.Getwd()
	if err != nil {
		logger.Fatal().Err(err).Msg("error retrieving working directory")
	}

	keysPath := path.Join(wd, ".secrets")
	logger = logger.With().
		Dict(
			"secrets_manager",
			zerolog.Dict().
				Int("backups", Secrets.Backups).
				Dur("update_interval", Secrets.UpdateInterval).
				Str("type", "local storage").
				Str("path", keysPath).
				Str("prefix", Secrets.Prefix),
		).
		Logger()

	return dao.NewFileSystemSecretKeysRepository(keysPath, Secrets.Prefix), logger
}
