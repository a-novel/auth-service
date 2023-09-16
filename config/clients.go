package config

import (
	"github.com/a-novel/go-api-clients"
	"github.com/rs/zerolog"
	"net/url"
)

func GetAuthorizationsClient(logger zerolog.Logger) apiclients.AuthorizationsClient {
	authorizationsURL, err := new(url.URL).Parse(API.External.AuthorizationsAPI)
	if err != nil {
		logger.Fatal().Err(err).Msg("could not parse auth API URL")
	}

	return apiclients.NewAuthorizationsClient(authorizationsURL)
}
