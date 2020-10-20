package repositories

import (
	"time"

	"github.com/google/uuid"
	"transport/lib/database"

	"telegram-bot/app/models"
)

type IMessageRepository interface {
	Create(msg *models.Message) error
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

func (a *messageRepo) Create(msg *models.Message) error {
	msg.ID = uuid.New().String()
	msg.CreatedTime = time.Now().Format(time.RFC3339Nano)
	msg.UpdatedTime = time.Now().Format(time.RFC3339Nano)

	err := a.db.InsertOne(models.CollectionMessage, &msg)
	if err != nil {
		return err
	}

	return nil
}
