package services

import (
	"context"

	"transport/lib/utils/paging"

	"telegram-bot/app/models"
	"telegram-bot/app/repositories"
)

type ActionService interface {
	Retrieve(ctx context.Context, name string) (*models.Action, error)
	List(ctx context.Context, name string) (*[]models.Action, *paging.Paging, error)
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
