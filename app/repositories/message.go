package repositories

import (
	"time"

	"github.com/google/uuid"
	"gopkg.in/mgo.v2/bson"
	"transport/lib/database"

	"telegram-bot/app/models"
	"telegram-bot/pkg/utils"
)

type IMessageRepository interface {
	Retrieve(id string) (*models.Message, error)
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

func (m *messageRepo) Retrieve(id string) (*models.Message, error) {
	msg := models.Message{}
	query := bson.M{"id": id}
	err := m.db.FindOne(models.CollectionMessage, query, DefaultSortField, &msg)
	if err != nil {
		return nil, utils.BOTGetDataNotfound.New()
	}

	return &msg, nil
}

func (m *messageRepo) Create(msg *models.Message) error {
	msg.ID = uuid.New().String()
	msg.CreatedTime = time.Now().Format(time.RFC3339Nano)
	msg.UpdatedTime = time.Now().Format(time.RFC3339Nano)

	err := m.db.InsertOne(models.CollectionMessage, &msg)
	if err != nil {
		return err
	}

	return nil
}
