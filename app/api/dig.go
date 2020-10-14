package api

import "go.uber.org/dig"

func Inject(container *dig.Container) error {
	_ = container.Provide(NewAction)
	_ = container.Provide(NewMessage)

	return nil
}
