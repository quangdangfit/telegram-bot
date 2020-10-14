package services

import (
	"context"

	"telegram-bot/app/models"
	"telegram-bot/app/repositories"
	"telegram-bot/app/schema"
)

type MessageService interface {
	Create(ctx context.Context, body *schema.MessageCreateParam) (*models.Message, error)
}

type messageService struct {
	messageRepo repositories.IMessageRepository
}

func NewMessageService(messageRepo repositories.IMessageRepository) MessageService {
	service := messageService{
		messageRepo: messageRepo,
	}
	return &service
}

func (m *messageService) Create(ctx context.Context, body *schema.MessageCreateParam) (*models.Message, error) {
	result, err := m.messageRepo.Create(body)
	if err != nil {
		return nil, err
	}

	return result, nil
}
