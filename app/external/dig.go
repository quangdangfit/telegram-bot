package external

import (
	"go.uber.org/dig"
	"transport/lib/thttp/httpclient"
)

func Inject(container *dig.Container) error {
	_ = container.Provide(httpclient.NewClient)
	_ = container.Provide(New)

	return nil
}
