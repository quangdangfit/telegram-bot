package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/dig"
	"transport/lib/ctx/ctxutils"
	"transport/lib/tconfig"
	"transport/lib/utils/logger"
)

func Initialize(container *dig.Container) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	app := gin.Default()
	app.Use(tconfig.CorsGinMiddleware())
	app.Use(ctxutils.GinGenContextMiddleware())

	err := RegisterAPI(app, container)
	if err != nil {
		logger.Error("Failed to register API: ", err)
	}

	if !viper.GetBool("publisher.ignore_store") &&
		viper.GetString("ts_service.env") != "production" {
		RegisterDocs(app)
	}

	return app
}

func InitializeCron(container *dig.Container) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	app := gin.Default()
	app.Use(tconfig.CorsGinMiddleware())
	app.Use(ctxutils.GinGenContextMiddleware())
	return app
}
