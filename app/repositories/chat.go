package repositories

import (
	"gopkg.in/mgo.v2/bson"
	"transport/lib/database"

	"telegram-bot/app/models"
	"telegram-bot/pkg/utils"
)

type IChatRepository interface {
	Retrieve(id int64) (*models.Chat, error)
	Create(chat *models.Chat) error
}

type chatRepo struct {
	db database.MongoDB
}

func NewChatRepository(db database.MongoDB) IChatRepository {
	ensureInIndex(db)
	return &chatRepo{db: db}
}

func (u *chatRepo) Retrieve(id int64) (*models.Chat, error) {
	chat := models.Chat{}
	query := bson.M{"id": id}
	err := u.db.FindOne(models.CollectionChat, query, DefaultSortField, &chat)
	if err != nil {
		return nil, utils.BOTGetDataNotfound.New()
	}

	return &chat, nil
}

func (u *chatRepo) Create(chat *models.Chat) error {
	err := u.db.InsertOne(models.CollectionChat, &chat)
	if err != nil {
		return utils.BOTAddDataError.New()
	}
	return nil
}
