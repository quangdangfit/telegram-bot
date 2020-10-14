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

	"transport/ems/app/models"
	"transport/ems/app/schema"
	"transport/ems/config"
	"transport/ems/pkg/utils"
)

type OutRepository interface {
	GetMessageByID(id string) (*models.OutMessage, error)
	GetMessage(query *schema.OutMessageQueryParam) (*models.OutMessage, error)
	GetMessages(query *schema.OutMessageQueryParam) (*[]models.OutMessage, *paging.Paging, error)
	GetMessagesInStatus(query *schema.OutMessageInStatusQueryParam) (*[]models.OutMessage, *paging.Paging, error)
	AddMessage(message *models.OutMessage) error
	UpdateMessage(id string, message *models.OutMessage) error
	DeleteMessage(id string) error
	ArchivedMessage(message *models.OutMessage) error
}

type outRepo struct {
	db database.MongoDB
}

func ensureOutIndex(db database.MongoDB) {
	indexOutId := mgo.Index{
		Key:        []string{"id"},
		Unique:     true,
		DropDups:   false,
		Background: true,
		Sparse:     true,
		Name:       "out_message_id_index",
	}
	db.DropIndex(models.CollectionOutMessage, indexOutId.Name)
	db.EnsureIndex(models.CollectionOutMessage, indexOutId)

	indexOutObj := mgo.Index{
		Key:        []string{"origin_code", "origin_model"},
		Unique:     false,
		DropDups:   false,
		Background: true,
		Sparse:     true,
		Name:       "out_origin_object_index",
	}
	db.DropIndex(models.CollectionOutMessage, indexOutObj.Name)
	db.EnsureIndex(models.CollectionOutMessage, indexOutObj)

	indexFind := mgo.Index{
		Key:        []string{"status", "created_time", "routing_key"},
		Unique:     false,
		DropDups:   false,
		Background: true,
		Sparse:     true,
		Name:       "out_status_index",
	}
	db.DropIndex(models.CollectionOutMessage, indexFind.Name)
	db.EnsureIndex(models.CollectionOutMessage, indexFind)
}

func NewOutRepository(db database.MongoDB) OutRepository {
	ensureOutIndex(db)
	return &outRepo{db: db}
}

func (o *outRepo) getCollection(createTime time.Time) string {
	return utils.GetArchivedCollection(models.CollectionOutMessage, createTime)
}

func (o *outRepo) GetMessageByID(id string) (*models.OutMessage, error) {
	message := models.OutMessage{}
	query := bson.M{"id": id}

	err := o.db.FindOne(models.CollectionOutMessage, query, "-_id", &message)
	if err != nil {
		return nil, utils.EMSGetDataError.New()
	}

	return &message, nil
}

func (o *outRepo) GetMessage(query *schema.OutMessageQueryParam) (*models.OutMessage, error) {
	message := models.OutMessage{}

	var mapQuery map[string]interface{}
	data, err := bson.Marshal(query)
	if err != nil {
		return nil, utils.EMSSerializeError.New()
	}
	bson.Unmarshal(data, &mapQuery)

	if query.Sort == "" {
		query.Sort = config.DefaultSortField
	}

	err = o.db.FindOne(models.CollectionOutMessage, mapQuery, query.Sort, &message)
	if err != nil {
		return nil, utils.EMSGetDataError.New()
	}

	return &message, nil
}

func (o *outRepo) GetMessages(query *schema.OutMessageQueryParam) (*[]models.OutMessage, *paging.Paging, error) {
	var message []models.OutMessage

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
		query.Sort = config.DefaultSortField
	}

	collection := models.CollectionOutMessage
	if query.CreatedTime != "" {
		createdTime, err := time.Parse("2006-01", query.CreatedTime)
		if err != nil {
			return nil, nil, utils.EMSParseTimeError.New()
		}
		collection = o.getCollection(createdTime)
	}

	pageInfo, err := o.db.FindManyPaging(collection, mapQuery, query.Sort, query.Page, query.Limit, &message)
	if err != nil {
		return nil, nil, utils.EMSGetDataError.New()
	}

	return &message, pageInfo, nil
}

func (o *outRepo) GetMessagesInStatus(query *schema.OutMessageInStatusQueryParam) (*[]models.OutMessage, *paging.Paging, error) {
	var message []models.OutMessage

	mapQuery := map[string]interface{}{}
	if query.Page <= 0 {
		query.Page = 1
	}

	if query.Limit <= 0 {
		query.Limit = config.DefaultPageSize
	}

	if query.Sort == "" {
		query.Sort = config.DefaultSortField
	}

	if query.Status != nil {
		mapQuery["status"] = bson.M{"$in": query.Status}
	}

	if query.CreatedDays > 0 {
		mapQuery["created_time"] = bson.M{"$lte": time.Now().UTC().AddDate(0, 0, -query.CreatedDays).Format(time.RFC3339Nano)}
	}

	pageInfo, err := o.db.FindManyPaging(models.CollectionOutMessage, mapQuery, query.Sort, query.Page, query.Limit, &message)
	if err != nil {
		return nil, nil, utils.EMSGetDataError.New()
	}

	return &message, pageInfo, nil
}

func (o *outRepo) AddMessage(message *models.OutMessage) error {
	message.CreatedTime = time.Now().UTC().Format(time.RFC3339Nano)
	message.UpdatedTime = time.Now().UTC().Format(time.RFC3339Nano)
	message.ID = uuid.New().String()

	var payload map[string]interface{}
	data, err := bson.Marshal(message)
	if err != nil {
		logger.Error("Marshal fail: ", err)
		return utils.EMSSerializeError.New()
	}
	bson.Unmarshal(data, &payload)

	err = o.db.InsertOne(models.CollectionOutMessage, payload)
	if err != nil {
		logger.Error("Cannot create out message: ", err)
		return utils.EMSAddDataError.New()
	}
	return nil
}

func (o *outRepo) UpdateMessage(id string, message *models.OutMessage) error {
	message.UpdatedTime = time.Now().UTC().Format(time.RFC3339Nano)

	var payload map[string]interface{}
	data, err := bson.Marshal(message)
	if err != nil {
		logger.Error("Marshal fail: ", err)
		return utils.EMSSerializeError.New()
	}
	bson.Unmarshal(data, &payload)

	change := bson.M{"$set": payload}
	query := bson.M{"id": id}
	err = o.db.UpdateOne(models.CollectionOutMessage, query, change)
	if err != nil {
		logger.Error("Cannot update out message: ", err)
		return utils.EMSUpdateDataError.New()
	}

	return nil
}

func (o *outRepo) DeleteMessage(id string) error {
	query := bson.M{"id": id}
	err := o.db.DeleteOne(models.CollectionOutMessage, query)
	if err != nil {
		logger.Error("Cannot delete out message: ", err)
		return utils.EMSDeleteDataError.New()
	}

	return nil
}

func (o *outRepo) ArchivedMessage(message *models.OutMessage) error {
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

	err = o.db.InsertOne(o.getCollection(createdTime), payload)
	if err != nil {
		logger.Error("Cannot add out message: ", err)
		return utils.EMSUpdateDataError.New()
	}

	return o.DeleteMessage(message.ID)
}
