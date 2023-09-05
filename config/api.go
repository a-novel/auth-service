package config

import (
	_ "embed"
	"github.com/a-novel/go-framework/middlewares"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"log"
)

//go:embed api-dev.yml
var apiDevFile []byte

//go:embed api-prod.yml
var apiProdFile []byte

type ApiConfig struct {
	Host         string `yaml:"host"`
	HostInternal string `yaml:"hostInternal"`
	Port         int    `yaml:"port"`
	PortInternal int    `yaml:"portInternal"`
}

var API *ApiConfig

func init() {
	cfg := new(ApiConfig)

	if err := loadEnv(EnvLoader{ProdENV: apiProdFile, DevENV: apiDevFile}, cfg); err != nil {
		log.Fatalf("error loading api configuration: %v\n", err)
	}

	API = cfg
}

func GetRouter(logger zerolog.Logger) *gin.Engine {
	router := gin.New()
	router.Use(gin.RecoveryWithWriter(logger), middlewares.Logger(logger, Deploy.ProjectID), cors.New(Cors))

	if ENV == ProdENV {
		gin.SetMode(gin.ReleaseMode)
		router.TrustedPlatform = gin.PlatformGoogleAppEngine
	}

	return router
}
