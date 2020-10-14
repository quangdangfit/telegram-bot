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

type OutMsg struct {
	service services.OutService
}

func NewOutMsg(service services.OutService) *OutMsg {
	return &OutMsg{service: service}
}

// Publish godoc
// @Tags Internal
// @Summary publish message to amqp
// @Description api publish out message to amqp
// @Accept  json
// @Produce json
// @Param Body body schema.OutMessageBodyParam true "Body"
// @Security BasicAuth
// @Success 200 {object} thttp.BaseResponse
// @Router /internal/messages [post]
func (o *OutMsg) Publish(c *gin.Context) thttp.Response {
	var body schema.OutMessageBodyParam
	if err := c.Bind(&body); err != nil {
		logger.Error("Publisher failed to bind body: ", err)
		return thttp.Response{
			Error: errors.BadRequest.New(),
		}
	}

	validate := validator.New()
	if err := validate.Validate(body); err != nil {
		logger.Error("Publisher body is invalid ", err)
		return thttp.Response{
			Error: errors.BadRequest.New(),
		}
	}

	err := o.service.Publish(c, &body)
	if err != nil {
		logger.Error("Publisher cannot publish message ", err)
		return thttp.Response{
			Error: err,
		}
	}

	return thttp.Response{
		Error: errors.Success.New(),
	}
}

// GetMessages godoc
// @Tags Out Messages
// @Summary api get list out messages
// @Description api get list out messages
// @Accept  json
// @Produce json
// @Param Query query schema.OutMessageQueryParam true "Query"
// @Security ApiKeyAuth
// @Success 200 {object} thttp.BaseResponse
// @Router /private/out_messages [get]
func (o *OutMsg) GetMessages(c *gin.Context) thttp.Response {
	var query schema.OutMessageQueryParam
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

	result, pageInfo, err := o.service.GetMessages(c, &query)
	var data []schema.OutMessage
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

// GetMessageByID godoc
// @Tags Out Messages
// @Summary api update out message by id
// @Description api update out message by id
// @Accept  json
// @Produce json
// @Param id path string true "Message ID"
// @Security ApiKeyAuth
// @Success 200 {object} thttp.BaseResponse
// @Router /private/out_messages/{id} [get]
func (o *OutMsg) GetMessageByID(c *gin.Context) thttp.Response {
	id := c.Param("id")
	if id == "" {
		logger.Error("Missing message id")
		return thttp.Response{
			Error: errors.BadRequest.New(),
		}
	}

	result, err := o.service.GetMessageByID(c, id)
	return thttp.Response{
		Error: err,
		Data:  result,
	}
}

// UpdateMessages godoc
// @Tags Out Messages
// @Summary api update out message
// @Description api update out message
// @Accept  json
// @Produce json
// @Param id path string true "Message ID"
// @Param body body schema.OutMessageBodyUpdateParam true "Body"
// @Security ApiKeyAuth
// @Success 200 {object} thttp.BaseResponse
// @Router /private/out_messages/{id} [put]
func (o *OutMsg) UpdateMessages(c *gin.Context) thttp.Response {
	id := c.Param("id")
	if id == "" {
		logger.Error("Missing message id")
		return thttp.Response{
			Error: errors.BadRequest.New(),
		}
	}

	var body schema.OutMessageBodyUpdateParam
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

	result, err := o.service.UpdateMessage(c, id, &body)
	return thttp.Response{
		Error: err,
		Data:  result,
	}
}
