package config

import (
	_ "embed"
	"log"
)

//go:embed mailer.yml
var mailerFile []byte

//go:embed mailer-dev.yml
var mailerDevFile []byte

type MailerConfig struct {
	APIKey  string `yaml:"apiKey"`
	Sandbox bool   `yaml:"sandbox"`
	Sender  struct {
		Email string `yaml:"email"`
		Name  string `yaml:"name"`
	} `yaml:"sender"`
	Templates struct {
		EmailValidation string `yaml:"emailValidation"`
		EmailUpdate     string `yaml:"emailUpdate"`
		PasswordReset   string `yaml:"passwordReset"`
	} `yaml:"templates"`
}

var Mailer *MailerConfig

func init() {
	cfg := new(MailerConfig)

	if err := loadEnv(EnvLoader{DefaultENV: mailerFile, DevENV: mailerDevFile}, cfg); err != nil {
		log.Fatalf("error loading app configuration: %v\n", err)
	}

	Mailer = cfg
}
