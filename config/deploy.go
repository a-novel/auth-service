package config

import (
	_ "embed"
	"log"
)

//go:embed deploy-prod.yml
var deployProdFile []byte

type DeployConfig struct {
	ProjectID string `yaml:"projectID"`
	Buckets   struct {
		SecretKeys string `yaml:"secretKeys"`
	} `yaml:"buckets"`
}

var Deploy *DeployConfig

func init() {
	cfg := new(DeployConfig)

	if err := loadEnv(EnvLoader{ProdENV: deployProdFile}, cfg); err != nil {
		log.Fatalf("error loading deploy configuration: %v\n", err)
	}

	Deploy = cfg
}
