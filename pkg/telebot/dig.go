package telebot

import "go.uber.org/dig"

func Inject(container *dig.Container) error {
	_ = container.Provide(NewTeleBot)

	return nil
}
