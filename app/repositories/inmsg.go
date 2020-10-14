package repositories

import (
	"gopkg.in/mgo.v2/bson"
	"transport/lib/database"
	"transport/lib/utils/paging"

	"telegram-bot/app/models"
)

const (
	DefaultSortField = "-created_time"
)

type IActionRepository interface {
	Retrieve(name string) (*models.Action, error)
	List(name string) (*[]models.Action, *paging.Paging, error)
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
	message := models.Action{}
	query := bson.M{"name": name}
	err := a.db.FindOne(models.CollectionAction, query, DefaultSortField, &message)
	if err != nil {
		return nil, err
	}

	return &message, nil
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
