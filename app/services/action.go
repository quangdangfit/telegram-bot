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

type actionService struct {
	actionRepo repositories.IActionRepository
}

func NewActionService(actionRepo repositories.IActionRepository) ActionService {
	service := actionService{
		actionRepo: actionRepo,
	}
	return &service
}

func (a *actionService) Retrieve(ctx context.Context, name string) (*models.Action, error) {
	result, err := a.actionRepo.Retrieve(name)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (a *actionService) List(ctx context.Context, name string) (*[]models.Action, *paging.Paging, error) {
	result, pageInfo, err := a.actionRepo.List(name)
	if err != nil {
		return nil, nil, err
	}

	return result, pageInfo, nil
}

func (a *actionService) Create(ctx context.Context, body *schema.ActionCreateParam) (*models.Action, error) {
	result, err := a.actionRepo.Create(body)
	if err != nil {
		return nil, err
	}

	return result, nil
}
