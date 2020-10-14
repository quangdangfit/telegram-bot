package api

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"transport/lib/errors"
	"transport/lib/thttp"
	"transport/lib/utils/logger"
	"transport/lib/validator"

	"telegram-bot/app/schema"
	"telegram-bot/app/services"
)

type Action struct {
	service services.ActionService
}

func NewAction(service services.ActionService) *Action {
	return &Action{service: service}
}

// List Actions godoc
// @Tags Actions
// @Summary api get list actions
// @Description api get list actions
// @Accept  json
// @Produce json
// @Param Query query schema.ActionQueryParam true "Query"
// @Security ApiKeyAuth
// @Success 200 {object} thttp.BaseResponse
// @Router /private/actions [get]
func (a *Action) List(c *gin.Context) thttp.Response {
	var query schema.ActionQueryParam
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

	result, pageInfo, err := a.service.List(c, query.Name)
	var data []schema.Action
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

// Create Action godoc
// @Tags Actions
// @Summary api create action
// @Description api create action
// @Accept  json
// @Produce json
// @Param Body body schema.ActionCreateParam true "Body"
// @Security BasicAuth
// @Success 200 {object} thttp.BaseResponse
// @Router /private/actions [post]
func (a *Action) Create(c *gin.Context) thttp.Response {
	var body schema.ActionCreateParam
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

	result, err := a.service.Create(c, &body)
	return thttp.Response{
		Error: err,
		Data:  result,
	}
}
