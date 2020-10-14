package services

import (
	"context"

	"telegram-bot/app/models"
	"telegram-bot/app/repositories"
	"telegram-bot/app/schema"
	"telegram-bot/pkg/telebot"
)

type MessageService interface {
	Create(ctx context.Context, body *schema.MessageCreateParam) (*models.Message, error)
}

type messageService struct {
	messageRepo repositories.IMessageRepository
	tele        telebot.TelegramBot
}

func NewMessageService(messageRepo repositories.IMessageRepository, tele telebot.TelegramBot) MessageService {
	service := messageService{
		messageRepo: messageRepo,
		tele:        tele,
	}
	return &service
}

func (m *messageService) Create(ctx context.Context, body *schema.MessageCreateParam) (*models.Message, error) {
	result, err := m.messageRepo.Create(body)
	if err != nil {
		return nil, err
	}

	m.tele.Send(ctx, nil)
	return result, nil
}
