package repositories

import (
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"gopkg.in/mgo.v2/bson"
	"transport/lib/database"
	"transport/lib/utils/paging"

	"telegram-bot/app/models"
	"telegram-bot/app/schema"
)

const (
	DefaultSortField = "-created_time"
)

type IActionRepository interface {
	Retrieve(name string) (*models.Action, error)
	List(name string) (*[]models.Action, *paging.Paging, error)
	Create(body *schema.ActionCreateParam) (*models.Action, error)
}

type actionRepo struct {
	db database.MongoDB
}

func ensureInIndex(db database.MongoDB) {
}

func NewActionRepository(db database.MongoDB) IActionRepository {
	ensureInIndex(db)
	return &actionRepo{db: db}
}

func (a *actionRepo) Retrieve(name string) (*models.Action, error) {
	action := models.Action{}
	query := bson.M{"name": name}
	err := a.db.FindOne(models.CollectionAction, query, DefaultSortField, &action)
	if err != nil {
		return nil, err
	}

	return &action, nil
}

func (a *actionRepo) List(name string) (*[]models.Action, *paging.Paging, error) {
	var actions []models.Action
	query := bson.M{"name": name}
	pageInfo, err := a.db.FindManyPaging(models.CollectionAction, query, DefaultSortField, 1, 25, &actions)
	if err != nil {
		return nil, nil, err
	}

	return &actions, pageInfo, nil
}

func (a *actionRepo) Create(body *schema.ActionCreateParam) (*models.Action, error) {
	action := models.Action{
		Model: models.Model{
			ID:          uuid.New().String(),
			CreatedTime: time.Now().UTC().Format(time.RFC3339Nano),
			UpdatedTime: time.Now().UTC().Format(time.RFC3339Nano),
		},
	}
	copier.Copy(&action, &body)

	err := a.db.InsertOne(models.CollectionAction, &action)
	if err != nil {
		return nil, err
	}

	return &action, nil
}
