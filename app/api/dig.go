package api

import "go.uber.org/dig"

func Inject(container *dig.Container) error {
	_ = container.Provide(NewCron)
	_ = container.Provide(NewRouting)
	_ = container.Provide(NewOutMsg)
	_ = container.Provide(NewInMsg)

	return nil
}
