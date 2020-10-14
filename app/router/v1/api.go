package v1

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/dig"

	"telegram-bot/app/api"
)

func RegisterAPI(r *gin.Engine, container *dig.Container) error {
	err := container.Invoke(func(
		inMsg *api.Action,
	) error {
		//route := r.Group(path)
		//basicAuthMiddleware := basic.BasicMiddleware(basic.WithValidator(basic.GetDefaultValidator()))
		//jwtSecretKey := viper.GetString("ts_encryption.jwt_secret")
		//jwtMiddleware := jwt.JWTMiddleware(jwt.WithKey(jwtSecretKey))
		//
		//// Internal API
		//internal := route.Group("/internal")
		//internal.Use(basicAuthMiddleware)
		//
		//internal.POST("messages", ginwrapper.Wrap(outMsg.Publish))
		//
		//// Private API
		//private := route.Group("/private")
		//private.Use(jwtMiddleware)
		//
		//private.GET("/routing_keys", ginwrapper.Wrap(routing.GetRoutingKeys))
		//private.POST("/routing_keys", ginwrapper.Wrap(routing.AddRoutingKey))
		//private.PUT("/routing_keys/:id", ginwrapper.Wrap(routing.UpdateRoutingKey))
		//private.DELETE("/routing_keys/:id", ginwrapper.Wrap(routing.DeleteRoutingKey))
		//
		//private.GET("/out_messages", ginwrapper.Wrap(outMsg.GetMessages))
		//private.GET("/out_messages/:id", ginwrapper.Wrap(outMsg.GetMessageByID))
		//private.PUT("/out_messages/:id", ginwrapper.Wrap(outMsg.UpdateMessages))
		//
		//private.GET("/in_messages", ginwrapper.Wrap(inMsg.GetMessages))
		//private.PUT("/in_messages/:id", ginwrapper.Wrap(inMsg.UpdateMessages))

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
