package config

import (
	_ "embed"
	"log"
	"time"
)

//go:embed tokens.yml
var tokensFile []byte

type TokensConfig struct {
	TTL        time.Duration `yaml:"ttl"`
	RenewDelta time.Duration `yaml:"renewDelta"`
}

var Tokens *TokensConfig

func init() {
	cfg := new(TokensConfig)

	if err := loadEnv(EnvLoader{DefaultENV: tokensFile}, cfg); err != nil {
		log.Fatalf("error loading tokens configuration: %v\n", err)
	}

	Tokens = cfg
}
