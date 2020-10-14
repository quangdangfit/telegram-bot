package app

import (
	"go.uber.org/dig"
	"transport/lib/utils/logger"

	"transport/ems/app/api"
	"transport/ems/app/dbs"
	"transport/ems/app/external"
	"transport/ems/app/queue"
	"transport/ems/app/repositories"
	"transport/ems/app/services"
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

	// Inject externals
	err = external.Inject(container)
	if err != nil {
		logger.Error("Failed to inject services", err)
	}

	// Inject queue
	err = queue.Inject(container)
	if err != nil {
		logger.Error("Failed to inject queue", err)
	}

	// Inject APIs
	err = api.Inject(container)
	if err != nil {
		logger.Error("Failed to inject APIs", err)
	}

	return container
}
