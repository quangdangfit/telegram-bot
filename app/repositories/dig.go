package repositories

import (
	"go.uber.org/dig"
)

func Inject(container *dig.Container) error {
	_ = container.Provide(NewActionRepository)
	_ = container.Provide(NewMessageRepository)
	_ = container.Provide(NewChatRepository)
	return nil
}
