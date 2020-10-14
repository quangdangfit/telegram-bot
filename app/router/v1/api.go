package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/dig"
	"transport/lib/auth/basic"
	"transport/lib/auth/jwt"
	"transport/lib/thttp/ginwrapper"

	"telegram-bot/app/api"
)

func RegisterAPI(r *gin.Engine, container *dig.Container) error {
	err := container.Invoke(func(
		action *api.Action,
	) error {
		route := r.Group("/")
		basicAuthMiddleware := basic.BasicMiddleware(basic.WithValidator(basic.GetDefaultValidator()))
		jwtSecretKey := viper.GetString("ts_encryption.jwt_secret")
		jwtMiddleware := jwt.JWTMiddleware(jwt.WithKey(jwtSecretKey))

		// Internal API
		internal := route.Group("/internal")
		internal.Use(basicAuthMiddleware)

		//internal.POST("messages", ginwrapper.Wrap(outMsg.Publish))

		// Private API
		private := route.Group("/private")
		private.Use(jwtMiddleware)

		private.GET("/actions", ginwrapper.Wrap(action.List))
		private.POST("/actions", ginwrapper.Wrap(action.Create))

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
