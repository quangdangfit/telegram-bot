package app

import (
	"go.uber.org/dig"
	"transport/lib/utils/logger"

	"telegram-bot/app/api"
	"telegram-bot/app/dbs"
	"telegram-bot/app/repositories"
	"telegram-bot/app/services"
)

func BuildContainer() *dig.Container {
	container := dig.New()

	// Inject database
	err := dbs.Inject(container)
	if err != nil {
		logger.Error("Failed to inject database", err)
	}

	// Inject repositories
	err = repositories.Inject(container)
	if err != nil {
		logger.Error("Failed to inject repositories", err)
	}

	// Inject services
	err = services.Inject(container)
	if err != nil {
		logger.Error("Failed to inject services", err)
	}

	// Inject APIs
	err = api.Inject(container)
	if err != nil {
		logger.Error("Failed to inject APIs", err)
	}

	return container
}
