package model

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"msgCenter/config"
	"msgCenter/db"
	"time"
)

var CallbackRequestLogTable = "callback_request_log"

type CallbackRequestLog struct {
	Id primitive.ObjectID `bson:"_id"`
	ProjectName string `bson:"project_name"`
	EventKey string `bson:"event_key"`
	QueueName string `bson:"queue_name"`
	RequestHost string `bson:"request_host"`
	RequestPath string `bson:"request_path"`
	RequestType string `bson:"request_type"`
	CallbackRequestIsJson bool `bson:"callback_request_is_json"`
	RequestData string `bson:"request_data"`
	RequestRes bool `bson:"request_res"`
	RequestError string `bson:"request_error"`
	RequestStatus int `bson:"request_status"`
	RequestResponse string `bson:"request_response"`
	CreatedAt time.Time `bson:"created_at"`
	CreatedAtStr string `bson:"created_at_str"`
}

//获取mongo生成的collection
func GetCallbackLogCollection() *mongo.Collection {
	return db.MongodbConn.Database(config.MongodbInc.MongoDB).Collection(CallbackRequestLogTable)
}

//插入一条数据
func InsertLog(logData *CallbackRequestLog) (insertId primitive.ObjectID) {
	collection := db.MongodbConn.Database(config.MongodbInc.MongoDB).Collection(CallbackRequestLogTable)
	insertRest, err := collection.InsertOne(context.TODO(), logData)
	if err != nil {
		return
	}
	insertId = insertRest.InsertedID.(primitive.ObjectID)
	return
}

//获取回调记录列表
func GetCallbackLogs(projectName string, eventKey string, startDate time.Time, endDate time.Time, page int) (list []CallbackRequestLog, total int)  {
	//查询条件
	filter := bson.M{}
	if len(projectName) > 0 {
		filter["project_name"] = projectName
	}
	if len(eventKey) > 0 {
		filter["event_key"] = eventKey
	}
	if startDate != (time.Time{}) && endDate != (time.Time{}){
		filter["created_at"] = bson.M{"$gte": startDate, "$lte": endDate}
	}
	
	//分页
	skip := (page - 1) * config.PageSize
	findOptions := &options.FindOptions{}
	findOptions.SetLimit(int64(config.PageSize))
	findOptions.SetSkip(int64(skip))
	findOptions.SetSort(bson.M{"_id": -1})
	
	//查询数据列表
	collection := GetCallbackLogCollection()
	res, decodeErr := collection.Find(context.TODO(), filter, findOptions)
	if decodeErr == nil && res != nil {
		_ = res.All(context.TODO(), &list)
	}
	
	//数据总数
	count, err := collection.CountDocuments(context.TODO(), filter)
	if err == nil {
		total = int(count)
	}
	
	return
}