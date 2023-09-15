package main

import (
	"crypto/ed25519"
	"fmt"
	"github.com/a-novel/auth-service/config"
	"github.com/a-novel/auth-service/pkg/handlers"
	"github.com/a-novel/auth-service/pkg/services"
	"github.com/a-novel/go-apis"
)

func keyGen() (ed25519.PrivateKey, error) {
	_, private, err := ed25519.GenerateKey(nil)
	return private, err
}

func main() {
	logger := config.GetInternalLogger()

	secretKeysDAO, logger := config.GetSecretsRepository(logger)

	generateTokenService := services.NewGenerateTokenService(secretKeysDAO, config.Tokens.TTL)
	getTokenService := services.NewGetTokenStatusService(secretKeysDAO)
	introspectTokenService := services.NewIntrospectTokenService(generateTokenService, getTokenService, config.Tokens.RenewDelta)
	rotateSecretKeysService := services.NewRotateSecretKeysService(secretKeysDAO, keyGen, config.Secrets.Backups)

	introspectTokenHandler := handlers.NewIntrospectTokenHandler(introspectTokenService)
	rotateSecretKeysHandler := handlers.NewRotateSecretKeysHandler(rotateSecretKeysService)

	router := apis.GetRouter(apis.RouterConfig{
		Logger:    logger,
		ProjectID: config.Deploy.ProjectID,
		Prod:      config.ENV == config.ProdENV,
	})

	router.GET("/auth", introspectTokenHandler.Handle)
	router.POST("/rotate-keys", rotateSecretKeysHandler.Handle)

	if err := router.Run(fmt.Sprintf(":%d", config.API.PortInternal)); err != nil {
		logger.Fatal().Err(err).Msg("a fatal error occurred while running the internal API, and the server had to shut down")
	}
}
