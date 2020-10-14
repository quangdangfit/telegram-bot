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
	"transport/ems/config"
)

type Routing struct {
	service services.RoutingService
}

func NewRouting(service services.RoutingService) *Routing {
	return &Routing{service: service}
}

// GetRoutingKeys godoc
// @Tags Routing Keys
// @Summary api get list routing keys
// @Description api get list routing keys
// @Accept  json
// @Produce json
// @Param Query query schema.RoutingKeyQueryParam true "Query"
// @Security ApiKeyAuth
// @Success 200 {object} thttp.BaseResponse
// @Router /private/routing_keys [get]
func (r *Routing) GetRoutingKeys(c *gin.Context) thttp.Response {
	var query schema.RoutingKeyQueryParam
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

	if query.Page <= 0 {
		query.Page = 1
	}
	if query.Limit <= 0 {
		query.Limit = config.DefaultPageSize
	}

	result, pageInfo, err := r.service.GetRoutingKeys(c, &query)
	var data []schema.RoutingKey
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

// AddRoutingKey godoc
// @Tags Routing Keys
// @Summary api create routing key
// @Description api create routing key
// @Accept  json
// @Produce json
// @Param Body body schema.RoutingKeyBodyCreateParam true "Body"
// @Security ApiKeyAuth
// @Success 200 {object} thttp.BaseResponse
// @Router /private/routing_keys [post]
func (r *Routing) AddRoutingKey(c *gin.Context) thttp.Response {
	var body schema.RoutingKeyBodyCreateParam
	if err := c.ShouldBind(&body); err != nil {
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

	result, err := r.service.AddRoutingKey(c, &body)
	return thttp.Response{
		Error: err,
		Data:  result,
	}
}

// UpdateRoutingKey godoc
// @Tags Routing Keys
// @Summary api update routing key
// @Description api update routing key
// @Accept  json
// @Produce json
// @Param id path string true "Routing Key ID"
// @Param body body schema.RoutingKeyBodyUpdateParam true "Body"
// @Security ApiKeyAuth
// @Success 200 {object} thttp.BaseResponse
// @Router /private/routing_keys/{id} [put]
func (r *Routing) UpdateRoutingKey(c *gin.Context) thttp.Response {
	id := c.Param("id")
	if id == "" {
		logger.Error("Missing routing key id")
		return thttp.Response{
			Error: errors.BadRequest.New(),
		}
	}

	var body schema.RoutingKeyBodyUpdateParam
	if err := c.ShouldBind(&body); err != nil {
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

	result, err := r.service.UpdateRoutingKey(c, id, &body)
	return thttp.Response{
		Error: err,
		Data:  result,
	}
}

// DeleteRoutingKey godoc
// @Tags Routing Keys
// @Summary api delete routing key
// @Description api delete routing key
// @Accept  json
// @Produce json
// @Param id path string true "Routing Key ID"
// @Security ApiKeyAuth
// @Success 200 {object} thttp.BaseResponse
// @Router /private/routing_keys/{id} [delete]
func (r *Routing) DeleteRoutingKey(c *gin.Context) thttp.Response {
	id := c.Param("id")
	if id == "" {
		logger.Error("Missing routing key id")
		return thttp.Response{
			Error: errors.BadRequest.New(),
		}
	}

	err := r.service.DeleteRoutingKey(c, id)
	return thttp.Response{
		Error: err,
	}
}
