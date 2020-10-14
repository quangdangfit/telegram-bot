package api

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"transport/lib/errors"
	"transport/lib/thttp"
	"transport/lib/utils/logger"
	"transport/lib/validator"

	"transport/ems/app/schema"
	"transport/ems/app/services"
)

type InMsg struct {
	service services.InService
}

func NewInMsg(service services.InService) *InMsg {
	return &InMsg{service: service}
}

// GetMessages godoc
// @Tags In Messages
// @Summary api get list in messages
// @Description api get list in messages
// @Accept  json
// @Produce json
// @Param Query query schema.InMessageQueryParam true "Query"
// @Security ApiKeyAuth
// @Success 200 {object} thttp.BaseResponse
// @Router /private/in_messages [get]
func (i *InMsg) GetMessages(c *gin.Context) thttp.Response {
	var query schema.InMessageQueryParam
	if err := c.ShouldBindQuery(&query); err != nil {
		logger.Error("Failed to bind query: ", err)
		return thttp.Response{
			Error: errors.BadRequest.New(),
		}
	}

	validate := validator.New()
	if err := validate.Validate(query); err != nil {
		logger.Error("Query is invalid ", err)
		return thttp.Response{
			Error: errors.BadRequest.New(),
		}
	}

	result, pageInfo, err := i.service.GetMessages(c, &query)
	var data []schema.InMessage
	copier.Copy(&data, &result)
	rs := schema.ResponsePagingResult{
		Paging: pageInfo,
		Data:   data,
	}

	return thttp.Response{
		Error: err,
		Data:  rs,
	}
}

// UpdateMessages godoc
// @Tags In Messages
// @Summary api update in message
// @Description api update in message
// @Accept  json
// @Produce json
// @Param id path string true "Message ID"
// @Param body body schema.InMessageBodyUpdateParam true "Body"
// @Security ApiKeyAuth
// @Success 200 {object} thttp.BaseResponse
// @Router /private/in_messages/{id} [put]
func (i *InMsg) UpdateMessages(c *gin.Context) thttp.Response {
	id := c.Param("id")
	if id == "" {
		logger.Error("Missing message id")
		return thttp.Response{
			Error: errors.BadRequest.New(),
		}
	}

	var body schema.InMessageBodyUpdateParam
	if err := c.ShouldBind(&body); err != nil {
		logger.Error("Failed to bind body: ", err)
		return thttp.Response{
			Error: errors.BadRequest.New(),
		}
	}

	validate := validator.New()
	if err := validate.Validate(body); err != nil {
		logger.Error("Query is invalid ", err)
		return thttp.Response{
			Error: errors.BadRequest.New(),
		}
	}

	result, err := i.service.UpdateMessage(c, id, &body)
	return thttp.Response{
		Error: err,
		Data:  result,
	}
}
