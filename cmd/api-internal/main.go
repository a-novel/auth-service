package main

import (
	"fmt"
	"github.com/a-novel/auth-service/config"
	"github.com/a-novel/auth-service/pkg/handlers"
	"github.com/a-novel/auth-service/pkg/services"
	"github.com/a-novel/go-framework/security"
)

func main() {
	logger := config.GetInternalLogger()

	secretKeysDAO, logger := config.GetSecretsRepository(logger)

	rotateSecretKeysService := services.NewRotateSecretKeysService(secretKeysDAO, security.JWKKeyGen, config.Secrets.Backups)

	pingHandler := handlers.NewPingHandler()
	rotateSecretKeysHandler := handlers.NewRotateSecretKeysHandler(rotateSecretKeysService)

	router := config.GetRouter(logger)

	router.GET("/ping", pingHandler.Handle)
	router.POST("/rotate-keys", rotateSecretKeysHandler.Handle)

	if err := router.Run(fmt.Sprintf(":%d", config.API.PortInternal)); err != nil {
		logger.Fatal().Err(err).Msg("a fatal error occurred while running the internal API, and the server had to shut down")
	}
}
