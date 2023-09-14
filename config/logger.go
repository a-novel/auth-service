package config

import (
	"github.com/rs/zerolog"
	"os"
)

func GetLogger() zerolog.Logger {
	logger := zerolog.New(os.Stdout).
		With().
		Dict("application", zerolog.Dict().Str("name", App.Name).Str("env", ENV)).
		Dict(
			"mailer", zerolog.Dict().
				Str("sender_email", Mailer.Sender.Email).
				Str("sender_name", Mailer.Sender.Name).
				Bool("sandbox", Mailer.Sandbox),
		).
		Logger()

	switch ENV {
	case ProdENV:
		logger = logger.With().Timestamp().Logger()
	default:
		logger = logger.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	return logger
}

func GetInternalLogger() zerolog.Logger {
	logger := zerolog.New(os.Stdout).
		With().
		Dict("application", zerolog.Dict().Str("name", App.Name+"-internal").Str("env", ENV)).
		Logger()

	switch ENV {
	case ProdENV:
		logger = logger.With().Timestamp().Logger()
	default:
		logger = logger.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	return logger
}
