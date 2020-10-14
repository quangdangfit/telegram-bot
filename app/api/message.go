package api

import (
	"github.com/gin-gonic/gin"
	"transport/lib/errors"
	"transport/lib/thttp"
	"transport/lib/utils/logger"
	"transport/lib/validator"

	"telegram-bot/app/schema"
	"telegram-bot/app/services"
)

type Message struct {
	service services.MessageService
}

func NewMessage(service services.MessageService) *Message {
	return &Message{service: service}
}

func (m *Message) Create(c *gin.Context) thttp.Response {
	var body schema.MessageCreateParam
	if err := c.Bind(&body); err != nil {
		logger.Error("Failed to bind body: ", err)
		return thttp.Response{
			Error: errors.BadRequest.New(),
		}
	}

	validate := validator.New()
	if err := validate.Validate(body); err != nil {
		logger.Error("Body is invalid ", err)
		return thttp.Response{
			Error: errors.BadRequest.New(),
		}
	}

	result, err := m.service.Create(c, &body)
	return thttp.Response{
		Error: err,
		Data:  result,
	}
}
