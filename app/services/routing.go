package services

import (
	"context"
	"encoding/json"
	"reflect"

	"github.com/jinzhu/copier"
	"transport/lib/errors"
	"transport/lib/utils/paging"

	"transport/ems/app/models"
	"transport/ems/app/repositories"
	"transport/ems/app/schema"
	"transport/ems/pkg/utils"
)

type RoutingService interface {
	GetRoutingKeys(ctx context.Context, query *schema.RoutingKeyQueryParam) (*[]models.RoutingKey, *paging.Paging, error)
	AddRoutingKey(ctx context.Context, body *schema.RoutingKeyBodyCreateParam) (*models.RoutingKey, error)
	UpdateRoutingKey(ctx context.Context, id string, body *schema.RoutingKeyBodyUpdateParam) (*models.RoutingKey, error)
	DeleteRoutingKey(ctx context.Context, id string) error
}

type routingService struct {
	repo repositories.RoutingRepository
}

func NewRoutingService(repo repositories.RoutingRepository) RoutingService {
	service := routingService{
		repo: repo,
	}

	return &service
}

func (r *routingService) GetRoutingKeys(ctx context.Context, query *schema.RoutingKeyQueryParam) (*[]models.RoutingKey, *paging.Paging, error) {
	result, pageInfo, err := r.repo.GetRoutingKeysPaging(query)
	if err != nil {
		return nil, nil, err
	}

	return result, pageInfo, errors.Success.New()
}

func (r *routingService) AddRoutingKey(ctx context.Context, body *schema.RoutingKeyBodyCreateParam) (*models.RoutingKey, error) {
	var routing models.RoutingKey
	copier.Copy(&routing, &body)

	err := r.repo.AddRoutingKey(&routing)

	if err != nil {
		return nil, err
	}

	return &routing, errors.Success.New()
}

func (r *routingService) UpdateRoutingKey(ctx context.Context, id string, body *schema.RoutingKeyBodyUpdateParam) (*models.RoutingKey, error) {
	routing, err := r.repo.GetRoutingKeyByID(id)
	if routing == nil {
		return nil, utils.EMSGetDataNotfound.New()
	}

	var update models.RoutingKey
	copier.Copy(&update, &routing)
	data, err := json.Marshal(body)
	if err != nil {
		return nil, utils.EMSSerializeError.New()
	}
	json.Unmarshal(data, &update)

	if reflect.DeepEqual(*routing, update) {
		return routing, errors.Success.New()
	}

	err = r.repo.UpdateRoutingKey(id, &update)
	if err != nil {
		return nil, err
	}

	return &update, errors.Success.New()
}

func (r *routingService) DeleteRoutingKey(ctx context.Context, id string) error {
	err := r.repo.DeleteRoutingKey(id)
	if err != nil {
		return err
	}

	return errors.Success.New()
}
