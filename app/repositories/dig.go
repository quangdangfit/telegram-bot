package repositories

import (
	"go.uber.org/dig"
)

func Inject(container *dig.Container) error {
	_ = container.Provide(NewInMessageRepository)
	_ = container.Provide(NewOutRepository)
	_ = container.Provide(NewRoutingRepository)
	return nil
}
