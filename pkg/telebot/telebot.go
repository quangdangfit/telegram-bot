package telebot

import (
	"context"
	"fmt"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/spf13/viper"
	"transport/lib/utils/logger"

	"telegram-bot/app/models"
	"telegram-bot/app/repositories"
	"telegram-bot/pkg/utils"
)

const (
	ShowMoreButton = "↓ Show more"
	ShowLessButton = "↑ Less"
	ConfirmButton  = "✓ Confirm"
	RejectButton   = "✗ Reject"
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
	userRepo repositories.IUserRepository
}

func NewTeleBot(userRepo repositories.IUserRepository) TelegramBot {
	token := viper.GetString("telegram.token")

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		logger.Panic(err)
	}

	tele := &telebot{
		bot:      bot,
		userRepo: userRepo,
	}

	return tele
}

func (t *telebot) generateMarkup(ctx context.Context, data interface{}, less bool) tgbotapi.InlineKeyboardMarkup {
	showType := tgbotapi.NewInlineKeyboardButtonData(ShowMoreButton, "more")
	if less {
		showType = tgbotapi.NewInlineKeyboardButtonData(ShowLessButton, "less")
	}

	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			showType,
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(ConfirmButton, "confirm"),
			tgbotapi.NewInlineKeyboardButtonData(RejectButton, "reject"),
		),
	)
}

func (t *telebot) handleMarkup(ctx context.Context, update tgbotapi.Update) {
	if update.CallbackQuery == nil {
		return
	}

	callback := update.CallbackQuery
	switch callback.Data {
	case "more":
		updateMarkup := t.generateMarkup(ctx, nil, true)
		data := utils.Jsonify(data)
		edit := tgbotapi.EditMessageTextConfig{
			BaseEdit: tgbotapi.BaseEdit{
				ChatID:      callback.Message.Chat.ID,
				MessageID:   callback.Message.MessageID,
				ReplyMarkup: &updateMarkup,
			},
			Text: callback.Message.Text + "\n" + data,
		}
		t.bot.Send(edit)
		break
	case "less":
		updateMarkup := t.generateMarkup(ctx, nil, false)
		edit := tgbotapi.EditMessageTextConfig{
			BaseEdit: tgbotapi.BaseEdit{
				ChatID:      callback.Message.Chat.ID,
				MessageID:   callback.Message.MessageID,
				ReplyMarkup: &updateMarkup,
			},
			Text: "Phiên bàn giao mới được tạo.",
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
			Text: fmt.Sprintf("Phiên bàn giao %s đã bị từ chối bởi %s.", data["name"], callback.From.String()),
		}
		t.bot.Send(edit)
		break
	}
}

func (t *telebot) Send(ctx context.Context, msg *models.Message) {
	numericKeyboard := t.generateMarkup(ctx, nil, false)
	for _, chatID := range msg.Action.ChatID {
		msg := tgbotapi.NewMessage(chatID, msg.GetContent())
		msg.ReplyMarkup = numericKeyboard
		t.bot.Send(msg)
	}
}

func (t *telebot) Start(ctx context.Context, update *tgbotapi.Update) {
	from := update.Message.From
	if from == nil {
		return
	}

	u, _ := t.userRepo.Retrieve(from.UserName)
	if u != nil {
		logger.Info("User is already existed")
		return
	}

	user := models.User{
		ChatID:   int64(from.ID),
		Username: from.UserName,
	}

	t.userRepo.Create(&user)
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
			msg.Text = "Xin chào, tôi là bot của hệ thống GHN Transportation. Chúc bạn một ngày tốt lành."
		case "close":
			msg.Text = "Closed"
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
		}

		if msg.Text != "" {
			_, err := t.bot.Send(msg)
			if err != nil {
				logger.Panic(err)
			}
		}
	}
}
