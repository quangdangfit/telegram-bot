package repositories

import (
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"transport/lib/database"

	"telegram-bot/app/models"
	"telegram-bot/app/schema"
)

type IMessageRepository interface {
	Create(body *schema.MessageCreateParam) (*models.Message, error)
}

type messageRepo struct {
	db database.MongoDB
}

func ensureMsgIndex(db database.MongoDB) {
}

func NewMessageRepository(db database.MongoDB) IMessageRepository {
	ensureMsgIndex(db)
	return &messageRepo{db: db}
}

func (a *messageRepo) Create(body *schema.MessageCreateParam) (*models.Message, error) {
	message := models.Message{
		Model: models.Model{
			ID:          uuid.New().String(),
			CreatedTime: time.Now().UTC().Format(time.RFC3339Nano),
			UpdatedTime: time.Now().UTC().Format(time.RFC3339Nano),
		},
	}
	copier.Copy(&message, &body)

	err := a.db.InsertOne(models.CollectionMessage, &message)
	if err != nil {
		return nil, err
	}

	return &message, nil
}
