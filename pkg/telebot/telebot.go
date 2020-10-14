package telebot

import (
	"context"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/spf13/viper"
	"transport/lib/utils/logger"

	"telegram-bot/app/models"
)

type TelegramBot interface {
	Send(ctx context.Context, message *models.Message)
}

type telebot struct {
	bot *tgbotapi.BotAPI
}

func NewTeleBot() TelegramBot {
	token := viper.GetString("telegram.token")

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		logger.Panic(err)
	}

	tele := &telebot{
		bot: bot,
	}

	return tele
}

func (t *telebot) Send(ctx context.Context, message *models.Message) {
	msg := tgbotapi.NewMessage(670391246, "Có chuyến mới chờ bạn xác nhận")
	t.bot.Send(msg)
}
