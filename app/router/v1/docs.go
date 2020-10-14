package v1

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"telegram-bot/docs"
)

func RegisterDocs(e *gin.Engine) {
	path := getPath()

	docs.SwaggerInfo.BasePath = path
	e.GET(path+"/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
