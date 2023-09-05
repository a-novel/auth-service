package config

import (
	_ "embed"
	"log"
)

//go:embed app.yml
var appFile []byte

//go:embed app-dev.yml
var appDevFile []byte

//go:embed app-prod.yml
var appProdFile []byte

type AppConfig struct {
	Name     string `yaml:"name"`
	Frontend struct {
		URLs   []string `yaml:"urls"`
		Routes struct {
			ValidateEmail    string `yaml:"validateEmail"`
			ValidateNewEmail string `yaml:"validateNewEmail"`
			ResetPassword    string `yaml:"resetPassword"`
		} `yaml:"routes"`
	} `yaml:"frontend"`
}

var App *AppConfig

func init() {
	cfg := new(AppConfig)

	if err := loadEnv(EnvLoader{DefaultENV: appFile, DevENV: appDevFile, ProdENV: appProdFile}, cfg); err != nil {
		log.Fatalf("error loading app configuration: %v\n", err)
	}

	App = cfg
}
