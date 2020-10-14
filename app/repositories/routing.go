package repositories

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"transport/lib/database"
	"transport/lib/utils/logger"
	"transport/lib/utils/paging"

	"telegram-bot/app/models"
	"telegram-bot/app/schema"
	"telegram-bot/config"
	"telegram-bot/pkg/utils"
)

type RoutingRepository interface {
	GetRoutingKeyByID(id string) (*models.RoutingKey, error)
	GetRoutingKey(query *schema.RoutingKeyQueryParam) (*models.RoutingKey, error)
	GetRoutingKeys(query *schema.RoutingKeyQueryParam) (*[]models.RoutingKey, error)
	GetRoutingKeysPaging(query *schema.RoutingKeyQueryParam) (*[]models.RoutingKey, *paging.Paging, error)
	AddRoutingKey(routing *models.RoutingKey) error
	UpdateRoutingKey(id string, routing *models.RoutingKey) error
	DeleteRoutingKey(id string) error
}

type routingRepo struct {
	db database.MongoDB
}

func ensureRoutingIndex(db database.MongoDB) {
	indexId := mgo.Index{
		Key:        []string{"id"},
		Unique:     true,
		DropDups:   false,
		Background: true,
		Sparse:     true,
		Name:       "routing_key_id_index",
	}
	db.DropIndex(models.CollectionRoutingKey, indexId.Name)
	db.EnsureIndex(models.CollectionRoutingKey, indexId)

	indexFind := mgo.Index{
		Key:        []string{"name", "value", "group"},
		Unique:     true,
		DropDups:   false,
		Background: true,
		Sparse:     true,
		Name:       "routing_key_find_index",
	}
	db.DropIndex(models.CollectionRoutingKey, indexFind.Name)
	db.EnsureIndex(models.CollectionRoutingKey, indexFind)
}

func NewRoutingRepository(db database.MongoDB) RoutingRepository {
	ensureRoutingIndex(db)
	r := routingRepo{db: db}
	return &r
}

func (r *routingRepo) GetRoutingKey(query *schema.RoutingKeyQueryParam) (*models.RoutingKey, error) {
	var routingKey models.RoutingKey

	var mapQuery map[string]interface{}
	data, err := json.Marshal(query)
	if err != nil {
		return nil, utils.EMSSerializeError.New()
	}
	json.Unmarshal(data, &mapQuery)

	err = r.db.FindOne(models.CollectionRoutingKey, mapQuery, "-_id", &routingKey)
	if err != nil {
		return nil, utils.EMSGetDataError.New()
	}
	return &routingKey, nil
}

func (r *routingRepo) GetRoutingKeys(query *schema.RoutingKeyQueryParam) (*[]models.RoutingKey, error) {
	var routingKeys []models.RoutingKey

	var mapQuery map[string]interface{}
	data, err := json.Marshal(query)
	if err != nil {
		return nil, utils.EMSSerializeError.New()
	}
	json.Unmarshal(data, &mapQuery)

	err = r.db.FindMany(models.CollectionRoutingKey, mapQuery, "-_id", &routingKeys)
	if err != nil {
		return nil, utils.EMSGetDataError.New()
	}
	return &routingKeys, nil
}

func (r *routingRepo) GetRoutingKeysPaging(query *schema.RoutingKeyQueryParam) (*[]models.RoutingKey, *paging.Paging, error) {
	var routingKeys []models.RoutingKey

	var mapQuery map[string]interface{}
	data, err := json.Marshal(query)
	if err != nil {
		return nil, nil, utils.EMSSerializeError.New()
	}
	json.Unmarshal(data, &mapQuery)

	pageInfo, err := r.db.FindManyPaging(models.CollectionRoutingKey, mapQuery, "-_id", query.Page, query.Limit, &routingKeys)
	if err != nil {
		return nil, nil, utils.EMSGetDataError.New()
	}
	return &routingKeys, pageInfo, nil
}

func (r *routingRepo) GetRoutingKeyByID(id string) (*models.RoutingKey, error) {
	var routingKey models.RoutingKey

	query := bson.M{"id": id}
	err := r.db.FindOne(models.CollectionRoutingKey, query, "-_id", &routingKey)
	if err != nil {
		return nil, utils.EMSGetDataError.New()
	}
	return &routingKey, nil
}

func (r *routingRepo) AddRoutingKey(routing *models.RoutingKey) error {
	routing.CreatedTime = time.Now().UTC().Format(time.RFC3339Nano)
	routing.UpdatedTime = time.Now().UTC().Format(time.RFC3339Nano)
	routing.ID = uuid.New().String()
	routing.Active = true

	if routing.RetryTimes <= 0 {
		routing.RetryTimes = config.DefaultMaxRetryTimes
	}

	var value map[string]interface{}
	data, err := json.Marshal(routing)
	if err != nil {
		logger.Error("Marshal fail: ", err)
		return utils.EMSSerializeError.New()
	}
	json.Unmarshal(data, &value)

	err = r.db.InsertOne(models.CollectionRoutingKey, value)
	if err != nil {
		logger.Error("Cannot create routing key: ", err)
		return utils.EMSAddDataError.New()
	}
	return nil
}

func (r *routingRepo) UpdateRoutingKey(id string, routing *models.RoutingKey) error {
	routing.UpdatedTime = time.Now().UTC().Format(time.RFC3339Nano)

	var value map[string]interface{}
	data, err := json.Marshal(routing)
	if err != nil {
		logger.Error("Marshal fail: ", err)
		return utils.EMSSerializeError.New()
	}
	json.Unmarshal(data, &value)

	query := bson.M{"id": id}
	err = r.db.UpdateOne(models.CollectionRoutingKey, query, value)
	if err != nil {
		logger.Error("Cannot update routing key: ", err)
		return utils.EMSUpdateDataError.New()
	}
	return nil
}

func (r *routingRepo) DeleteRoutingKey(id string) error {
	query := bson.M{"id": id}
	err := r.db.DeleteOne(models.CollectionRoutingKey, query)
	if err != nil {
		logger.Error("Cannot delete routing key: ", err)
		return utils.EMSDeleteDataError.New()
	}
	return nil
}
