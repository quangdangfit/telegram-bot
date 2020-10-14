package services

import (
	"context"

	"transport/lib/utils/paging"

	"telegram-bot/app/models"
	"telegram-bot/app/repositories"
	"telegram-bot/app/schema"
)

type ActionService interface {
	Retrieve(ctx context.Context, name string) (*models.Action, error)
	List(ctx context.Context, name string) (*[]models.Action, *paging.Paging, error)
	Create(ctx context.Context, body *schema.ActionCreateParam) (*models.Action, error)
}

type inService struct {
	actionRepo repositories.IActionRepository
}

func NewActionService(actionRepo repositories.IActionRepository) ActionService {
	service := inService{
		actionRepo: actionRepo,
	}
	return &service
}

func (i *inService) Retrieve(ctx context.Context, name string) (*models.Action, error) {
	result, err := i.actionRepo.Retrieve(name)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (i *inService) List(ctx context.Context, name string) (*[]models.Action, *paging.Paging, error) {
	result, pageInfo, err := i.actionRepo.List(name)
	if err != nil {
		return nil, nil, err
	}

	return result, pageInfo, nil
}

func (i *inService) Create(ctx context.Context, body *schema.ActionCreateParam) (*models.Action, error) {
	result, err := i.actionRepo.Create(body)
	if err != nil {
		return nil, err
	}

	return result, nil
}
