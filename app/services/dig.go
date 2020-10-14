package services

import (
	"go.uber.org/dig"
)

func Inject(container *dig.Container) error {
	_ = container.Provide(NewActionService)
	_ = container.Provide(NewMessageService)

	return nil
}
