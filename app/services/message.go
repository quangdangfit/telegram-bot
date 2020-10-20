package services

import (
	"context"

	"github.com/jinzhu/copier"
	"transport/lib/errors"

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
	actionRepo  repositories.IActionRepository
	tele        telebot.TelegramBot
}

func NewMessageService(messageRepo repositories.IMessageRepository, actionRepo repositories.IActionRepository, tele telebot.TelegramBot) MessageService {
	service := messageService{
		messageRepo: messageRepo,
		actionRepo:  actionRepo,
		tele:        tele,
	}
	return &service
}

func (m *messageService) Create(ctx context.Context, body *schema.MessageCreateParam) (*models.Message, error) {
	var msg models.Message
	copier.Copy(&msg, &body)
	action, err := m.actionRepo.Retrieve(body.Action)
	if err != nil {
		return nil, err
	}

	msg.Action = *action
	err = m.messageRepo.Create(&msg)
	if err != nil {
		return nil, err
	}

	m.tele.Send(ctx, &msg)
	return &msg, errors.Success.New()
}
