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
func (i *Action) List(c *gin.Context) thttp.Response {
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

	result, pageInfo, err := i.service.List(c, query.Name)
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
