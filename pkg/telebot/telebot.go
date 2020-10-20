package telebot

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/spf13/viper"
	"transport/lib/utils/logger"

	"telegram-bot/app/models"
	"telegram-bot/app/repositories"
)

const (
	ConfirmButton = "✓ Confirm"
	RejectButton  = "✗ Reject"
	CancelButton  = "㊀ Cancel"
)

var data = map[string]interface{}{
	"name":       " 12ASDASD ASD2",
	"hub":        "HCM",
	"tag":        "Liên tỉnh",
	"cycle_type": "day",
	"cycle_time": 300,
	"start_time": "2020-08-09T09:45:51+07:00",
	"end_time":   "2020-09-13T09:45:51+07:00",
	"stoppoints": []map[string]interface{}{
		{
			"code":            "55",
			"name":            "Kho giao nhận Yên Minh_Hà Giang",
			"type":            "express",
			"address":         "HCM",
			"sort":            1,
			"est_distance":    0,
			"est_duration":    1322312,
			"est_day":         0,
			"est_time_in_at":  0,
			"est_time_out_at": 120,
		},
		{
			"code":            "55",
			"name":            "Kho giao nhận Yên Minh_Hà Giang",
			"type":            "express",
			"address":         "HCM",
			"sort":            2,
			"est_distance":    0,
			"est_duration":    1322312,
			"est_day":         1,
			"est_time_in_at":  0,
			"est_time_out_at": 120,
		}},
	"drivers": map[string]interface{}{
		"id":       5,
		"fullname": "Nguyễn Phát Lợi",
		"phone":    "983214701",
	},
	"vehicle": map[string]interface{}{
		"id":        2,
		"id_number": "51D-03913",
		"payload":   2300,
	},
}

type TelegramBot interface {
	Send(ctx context.Context, message *models.Message)
	Listen(ctx context.Context)
}

type telebot struct {
	bot      *tgbotapi.BotAPI
	chatRepo repositories.IChatRepository
	msgRepo  repositories.IMessageRepository
}

func NewTeleBot(chatRepo repositories.IChatRepository, msgRepo repositories.IMessageRepository) TelegramBot {
	token := viper.GetString("telegram.token")
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		logger.Panic(err)
	}

	tele := &telebot{
		bot:      bot,
		chatRepo: chatRepo,
		msgRepo:  msgRepo,
	}

	return tele
}

func (t *telebot) generateMarkupData(action string, id string) string {
	data := map[string]interface{}{
		"id":     id,
		"action": action,
	}
	b, _ := json.Marshal(data)
	return string(b)
}

func (t *telebot) generateMarkup(ctx context.Context, id string) *tgbotapi.InlineKeyboardMarkup {
	markup := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(ConfirmButton, t.generateMarkupData("confirm", id)),
			tgbotapi.NewInlineKeyboardButtonData(RejectButton, t.generateMarkupData("reject", id)),
			tgbotapi.NewInlineKeyboardButtonData(CancelButton, t.generateMarkupData("cancel", id)),
		),
	)
	return &markup
}

func (t *telebot) handleMarkup(ctx context.Context, update tgbotapi.Update) {
	if update.CallbackQuery == nil {
		return
	}

	callback := update.CallbackQuery
	strData := callback.Data
	var mapData map[string]string
	json.Unmarshal([]byte(strData), &mapData)

	if mapData["action"] == "" || mapData["id"] == "" {
		return
	}

	msg, err := t.msgRepo.Retrieve(mapData["id"])
	if err != nil {
		return
	}
	logger.Info(*msg)

	switch mapData["action"] {
	case "confirm":
		edit := tgbotapi.EditMessageTextConfig{
			BaseEdit: tgbotapi.BaseEdit{
				ChatID:      callback.Message.Chat.ID,
				MessageID:   callback.Message.MessageID,
				ReplyMarkup: nil,
			},
			Text: fmt.Sprintf("Phiên bàn giao %s đã được xác nhận bởi @%s.", data["name"], callback.From.String()),
		}
		t.bot.Send(edit)
		break
	case "reject":
		edit := tgbotapi.EditMessageTextConfig{
			BaseEdit: tgbotapi.BaseEdit{
				ChatID:      callback.Message.Chat.ID,
				MessageID:   callback.Message.MessageID,
				ReplyMarkup: nil,
			},
			Text: fmt.Sprintf("Phiên bàn giao %s đã bị từ chối bởi @%s.", data["name"], callback.From.String()),
		}
		t.bot.Send(edit)
		break
	case "cancel":
		edit := tgbotapi.EditMessageTextConfig{
			BaseEdit: tgbotapi.BaseEdit{
				ChatID:      callback.Message.Chat.ID,
				MessageID:   callback.Message.MessageID,
				ReplyMarkup: nil,
			},
			Text: fmt.Sprintf("Phiên bàn giao %s đã bị hủy bởi @%s.", data["name"], callback.From.String()),
		}
		t.bot.Send(edit)
		break
	}
}

func (t *telebot) Send(ctx context.Context, msg *models.Message) {
	numericKeyboard := t.generateMarkup(ctx, msg.ID)
	for _, chatID := range msg.Action.ChatID {
		msg := tgbotapi.NewMessage(chatID, msg.GetFullContent())
		msg.ReplyMarkup = numericKeyboard
		_, err := t.bot.Send(msg)

		if err != nil {
			logger.Error("Cannot send message ", err)
		}
	}
}

func (t *telebot) Start(ctx context.Context, update *tgbotapi.Update) {
	chat := update.Message.Chat
	if chat == nil {
		return
	}

	u, _ := t.chatRepo.Retrieve(chat.ID)
	if u != nil {
		logger.Info("Chat is already existed")
		return
	}

	c := models.Chat{
		ID:       chat.ID,
		Username: chat.UserName,
	}

	t.chatRepo.Create(&c)
}

func (t *telebot) Listen(ctx context.Context) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, _ := t.bot.GetUpdatesChan(u)

	for update := range updates {
		if update.CallbackQuery != nil {
			t.handleMarkup(ctx, update)
			continue
		}

		if update.Message == nil { // ignore non-Message updates
			continue
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)

		switch update.Message.Text {
		case "/start", "start":
			t.Start(ctx, &update)
			msg.Text = "Hệ thống đã ghi nhận tài khoản của bạn. Xin cám ơn."
		default:
			msg.Text = "Xin lỗi, hệ thống không hiểu lệnh của bạn."
		}

		if msg.Text != "" {
			_, err := t.bot.Send(msg)
			if err != nil {
				logger.Panic(err)
			}
		}
	}
}
