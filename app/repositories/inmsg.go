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

const (
	DefaultSortField = "-publish_time"
)

type InRepository interface {
	GetMessageByID(id string) (*models.InMessage, error)
	GetMessage(query *schema.InMessageQueryParam) (*models.InMessage, error)
	GetPreviousMessage(query *schema.InMessagePreviousQueryParam) (*models.InMessage, error)
	GetMessages(query *schema.InMessageQueryParam) (*[]models.InMessage, *paging.Paging, error)
	GetMessagesInStatus(query *schema.InMessageInStatusQueryParam) (*[]models.InMessage, *paging.Paging, error)
	AddMessage(message *models.InMessage) error
	UpdateMessage(message *models.InMessage) error
	DeleteMessage(id string) error
	ArchivedMessage(message *models.InMessage) error
}

type inRepo struct {
	db database.MongoDB
}

func ensureInIndex(db database.MongoDB) {
	indexInId := mgo.Index{
		Key:        []string{"id"},
		Unique:     true,
		DropDups:   false,
		Background: true,
		Sparse:     true,
		Name:       "in_message_id_index",
	}
	db.DropIndex(models.CollectionInMessage, indexInId.Name)
	db.EnsureIndex(models.CollectionInMessage, indexInId)

	indexInObj := mgo.Index{
		Key:        []string{"origin_code", "origin_model"},
		Unique:     false,
		DropDups:   false,
		Background: true,
		Sparse:     true,
		Name:       "in_origin_object_index",
	}
	db.DropIndex(models.CollectionInMessage, indexInObj.Name)
	db.EnsureIndex(models.CollectionInMessage, indexInObj)

	indexFind := mgo.Index{
		Key:        []string{"status", "created_time", "routing.name", "routing.value", "routing.group"},
		Unique:     false,
		DropDups:   false,
		Background: true,
		Sparse:     true,
		Name:       "in_status_index",
	}
	db.DropIndex(models.CollectionInMessage, indexFind.Name)
	db.EnsureIndex(models.CollectionInMessage, indexFind)
}

func NewInMessageRepository(db database.MongoDB) InRepository {
	ensureInIndex(db)
	return &inRepo{db: db}
}

func (i *inRepo) getCollection(createTime time.Time) string {
	return utils.GetArchivedCollection(models.CollectionInMessage, createTime)
}

func (i *inRepo) GetMessageByID(id string) (*models.InMessage, error) {
	message := models.InMessage{}
	query := bson.M{"id": id}
	err := i.db.FindOne(models.CollectionInMessage, query, "-_id", &message)
	if err != nil {
		return nil, utils.EMSGetDataError.New()
	}

	return &message, nil
}

func (i *inRepo) GetMessage(query *schema.InMessageQueryParam) (*models.InMessage, error) {
	message := models.InMessage{}

	var mapQuery map[string]interface{}
	data, err := json.Marshal(query)
	if err != nil {
		return nil, utils.EMSSerializeError.New()
	}
	json.Unmarshal(data, &mapQuery)

	if query.Sort == "" {
		query.Sort = DefaultSortField
	}

	err = i.db.FindOne(models.CollectionInMessage, mapQuery, query.Sort, &message)
	if err != nil {
		return nil, utils.EMSGetDataError.New()
	}

	return &message, nil
}

func (i *inRepo) GetPreviousMessage(query *schema.InMessagePreviousQueryParam) (*models.InMessage, error) {
	message := models.InMessage{}

	var mapQuery map[string]interface{}
	data, err := json.Marshal(query)
	if err != nil {
		return nil, utils.EMSSerializeError.New()
	}
	json.Unmarshal(data, &mapQuery)

	if query.Sort == "" {
		query.Sort = DefaultSortField
	}

	if query.PublishTime != "" {
		mapQuery["publish_time"] = bson.M{"$lt": query.PublishTime}
	}

	err = i.db.FindOne(models.CollectionInMessage, mapQuery, query.Sort, &message)
	if err != nil {
		return nil, utils.EMSGetDataError.New()
	}

	return &message, nil
}

func (i *inRepo) GetMessages(query *schema.InMessageQueryParam) (*[]models.InMessage, *paging.Paging, error) {
	var message []models.InMessage

	var mapQuery map[string]interface{}
	data, err := json.Marshal(query)
	if err != nil {
		return nil, nil, utils.EMSSerializeError.New()
	}
	json.Unmarshal(data, &mapQuery)

	if query.Page <= 0 {
		query.Page = 1
	}

	if query.Limit <= 0 {
		query.Limit = config.DefaultPageSize
	}

	if query.Sort == "" {
		query.Sort = DefaultSortField
	}

	collection := models.CollectionInMessage
	if query.CreatedTime != "" {
		createdTime, err := time.Parse("2006-01", query.CreatedTime)
		if err != nil {
			return nil, nil, utils.EMSParseTimeError.New()
		}
		collection = i.getCollection(createdTime)
	}

	pageInfo, err := i.db.FindManyPaging(collection, mapQuery, query.Sort, query.Page, query.Limit, &message)
	if err != nil {
		return nil, nil, utils.EMSGetDataError.New()
	}

	return &message, pageInfo, nil
}

func (i *inRepo) GetMessagesInStatus(query *schema.InMessageInStatusQueryParam) (*[]models.InMessage, *paging.Paging, error) {
	var message []models.InMessage

	mapQuery := map[string]interface{}{}
	if query.Page <= 0 {
		query.Page = 1
	}

	if query.Limit <= 0 {
		query.Limit = config.DefaultPageSize
	}

	if query.Sort == "" {
		query.Sort = DefaultSortField
	}

	if query.Status != nil {
		mapQuery["status"] = bson.M{"$in": query.Status}
	}

	if query.CreatedDays > 0 {
		mapQuery["created_time"] = bson.M{"$lte": time.Now().UTC().AddDate(0, 0, -query.CreatedDays).Format(time.RFC3339Nano)}
	}

	if query.UpdatedDays > 0 {
		mapQuery["updated_time"] = bson.M{"$lte": time.Now().UTC().AddDate(0, 0, -query.UpdatedDays).Format(time.RFC3339Nano)}
	}

	pageInfo, err := i.db.FindManyPaging(models.CollectionInMessage, mapQuery, query.Sort, query.Page, query.Limit, &message)
	if err != nil {
		return nil, nil, utils.EMSGetDataError.New()
	}

	return &message, pageInfo, nil
}

func (i *inRepo) AddMessage(message *models.InMessage) error {
	message.CreatedTime = time.Now().UTC().Format(time.RFC3339Nano)
	message.UpdatedTime = time.Now().UTC().Format(time.RFC3339Nano)
	message.ID = uuid.New().String()
	message.Attempts = 0

	var value map[string]interface{}
	data, err := json.Marshal(message)
	if err != nil {
		logger.Error("Marshal fail: ", err)
		return utils.EMSSerializeError.New()
	}
	json.Unmarshal(data, &value)

	err = i.db.InsertOne(models.CollectionInMessage, value)
	if err != nil {
		logger.Error("Cannot create in message: ", err)
		return utils.EMSAddDataError.New()
	}
	return nil
}

func (i *inRepo) UpdateMessage(message *models.InMessage) error {
	query := bson.M{"id": message.ID}

	var payload map[string]interface{}
	message.UpdatedTime = time.Now().UTC().Format(time.RFC3339Nano)

	data, err := json.Marshal(message)
	if err != nil {
		logger.Error("Marshal fail: ", err)
		return utils.EMSSerializeError.New()
	}
	json.Unmarshal(data, &payload)

	change := bson.M{"$set": payload}
	err = i.db.UpdateOne(models.CollectionInMessage, query, change)
	if err != nil {
		logger.Error("Cannot update in message: ", err)
		return utils.EMSUpdateDataError.New()
	}

	return nil
}

func (i *inRepo) DeleteMessage(id string) error {
	query := bson.M{"id": id}
	err := i.db.DeleteOne(models.CollectionInMessage, query)
	if err != nil {
		logger.Error("Cannot delete in message: ", err)
		return utils.EMSDeleteDataError.New()
	}

	return nil
}

func (i *inRepo) ArchivedMessage(message *models.InMessage) error {
	var payload map[string]interface{}
	data, err := bson.Marshal(message)
	if err != nil {
		logger.Error("Marshal fail: ", err)
		return utils.EMSSerializeError.New()
	}
	bson.Unmarshal(data, &payload)

	createdTime, err := time.Parse(time.RFC3339Nano, message.CreatedTime)
	if err != nil {
		return utils.EMSParseTimeError.New()
	}

	err = i.db.InsertOne(i.getCollection(createdTime), payload)
	if err != nil {
		logger.Error("Cannot add in message: ", err)
		return utils.EMSUpdateDataError.New()
	}

	return i.DeleteMessage(message.ID)
}
