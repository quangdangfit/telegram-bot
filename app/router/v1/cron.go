package v1

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/dig"
	"transport/lib/auth/basic"
	"transport/lib/thttp/ginwrapper"

	"transport/ems/app/api"
)

func RegisterCron(r *gin.Engine, container *dig.Container) error {
	path := getPath()
	err := container.Invoke(func(
		cron *api.Cron,
	) error {
		route := r.Group(path)
		basicAuthMiddleware := basic.BasicMiddleware(basic.WithValidator(basic.GetDefaultValidator()))

		cronRoute := route.Group("/cron")
		cronRoute.Use(basicAuthMiddleware)

		cronRoute.POST("/retry_previous", ginwrapper.Wrap(cron.RetryPrevious))
		cronRoute.POST("/resend", ginwrapper.Wrap(cron.Resend))
		cronRoute.POST("/retry", ginwrapper.Wrap(cron.Retry))
		cronRoute.POST("/archived", ginwrapper.Wrap(cron.ArchivedMessages))

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
