package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

var ENV = os.Getenv("ENV")

const (
	DevENV  = "dev"
	ProdENV = "prod"

	// DefaultENV is used to target all environments at once. This value should not be used as the actual content
	// of the ENV variable.
	DefaultENV = "all"
)

type EnvLoader map[string][]byte

func init() {
	if ENV == "" {
		ENV = DevENV
	}
}

func loadEnv(files EnvLoader, out interface{}) error {
	if ENV == "" {
		ENV = DevENV
	}

	for env, file := range files {
		if env != DefaultENV && env != ENV {
			continue
		}

		if err := yaml.Unmarshal([]byte(os.ExpandEnv(string(file))), out); err != nil {
			return err
		}
	}

	return nil
}
