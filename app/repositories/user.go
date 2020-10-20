package repositories

import (
	"time"

	"github.com/google/uuid"
	"gopkg.in/mgo.v2/bson"
	"transport/lib/database"

	"telegram-bot/app/models"
	"telegram-bot/pkg/utils"
)

type IUserRepository interface {
	Retrieve(chatID int64) (*models.User, error)
	Create(user *models.User) error
}

type userRepo struct {
	db database.MongoDB
}

func NewUserRepository(db database.MongoDB) IUserRepository {
	ensureInIndex(db)
	return &userRepo{db: db}
}

func (u *userRepo) Retrieve(chatID int64) (*models.User, error) {
	user := models.User{}
	query := bson.M{"chat_id": chatID}
	err := u.db.FindOne(models.CollectionUser, query, DefaultSortField, &user)
	if err != nil {
		return nil, utils.BOTGetDataNotfound.New()
	}

	return &user, nil
}

func (u *userRepo) Create(user *models.User) error {
	user.ID = uuid.New().String()
	user.CreatedTime = time.Now().UTC().Format(time.RFC3339Nano)
	user.UpdatedTime = time.Now().UTC().Format(time.RFC3339Nano)

	err := u.db.InsertOne(models.CollectionUser, &user)
	if err != nil {
		return utils.BOTAddDataError.New()
	}
	return nil
}
