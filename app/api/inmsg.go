package api

import (
	"telegram-bot/app/services"
)

type Action struct {
	service services.ActionService
}

func NewAction(service services.ActionService) *Action {
	return &Action{service: service}
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
//func (i *Action) GetMessages(c *gin.Context) thttp.Response {
//	var query schema.InMessageQueryParam
//	if err := c.ShouldBindQuery(&query); err != nil {
//		logger.Error("Failed to bind query: ", err)
//		return thttp.Response{
//			Error: errors.BadRequest.New(),
//		}
//	}
//
//	validate := validator.New()
//	if err := validate.Validate(query); err != nil {
//		logger.Error("Query is invalid ", err)
//		return thttp.Response{
//			Error: errors.BadRequest.New(),
//		}
//	}
//
//	result, pageInfo, err := i.service.GetMessages(c, &query)
//	var data []schema.InMessage
//	copier.Copy(&data, &result)
//	rs := schema.ResponsePagingResult{
//		Paging: pageInfo,
//		Data:   data,
//	}
//
//	return thttp.Response{
//		Error: err,
//		Data:  rs,
//	}
//}
