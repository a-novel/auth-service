package config

import (
	"github.com/a-novel/go-apis/clients"
	"github.com/rs/zerolog"
	"net/url"
)

func GetPermissionsClient(logger zerolog.Logger) apiclients.PermissionsClient {
	permissionsURL, err := new(url.URL).Parse(API.External.PermissionsAPI)
	if err != nil {
		logger.Fatal().Err(err).Msg("could not parse auth API URL")
	}

	return apiclients.NewPermissionsClient(permissionsURL)
}
