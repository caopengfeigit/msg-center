package model

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"msgCenter/config"
	"msgCenter/db"
)

var projectEventTable = "project_event"

type ProjectEvent struct {
	Id primitive.ObjectID `bson:"_id"`
	ProjectId string `bson:"project_id"`
	Name string `bson:"name"`
	Type string `bson:"type"`
	ExchangeType string `bson:"exchange_type"`
}

//获取mongo生成的collection
func GetProjectEventCollection() *mongo.Collection {
	return db.MongodbConn.Database(config.MongodbInc.MongoDB).Collection(projectEventTable)
}

//通过项目id和name获取一个消息类型配置
func GetOneProjectEventByIdName(projectId primitive.ObjectID, name string) (projectEvent ProjectEvent, code int) {
	collection := GetProjectEventCollection()
	filter := bson.D{{"project_id", projectId}, {"name", name}}
	decodeErr := collection.FindOne(context.TODO(), filter).Decode(&projectEvent)
	if decodeErr != nil {
		if projectEvent == (ProjectEvent{}) {
			code = 30002
		} else {
			code = 30004
		}
	}
	return
}

//通过类型获取一个消息类型配置
/*func GetOneProjectEventByType(typeStr string) (results []ProjectEvent) {
	collection := GetProjectEventCollection()
	filter := bson.D{{"type", typeStr}}
	cursor, _ := collection.Find(context.TODO(), filter)
	for cursor.Next(context.TODO()) {
		var item ProjectEvent
		err := cursor.Decode(&item)
		if err != nil {
			continue
		}
		results = append(results, item)
	}
	return
}*/

//获取PublishSubscribe队列数据的查询
func getPipeline(eventType string, eventId primitive.ObjectID) (pipeline mongo.Pipeline) {
	eventTypeConfigTable := getEventTypeConfigTable(eventType)
	project := getOutFields(eventType)
	var match bson.E
	if eventId.IsZero() {
		match = bson.E{"type", eventType}
	} else {
		match = bson.E{"_id", eventId}
	}
	pipeline = mongo.Pipeline {
		//match条件
		{
			{
				"$match" , bson.D {match},
			},
		},
		//第一个连表是查项目名，目的是生成唯一的队列名
		{
			{
				"$lookup", bson.D{
				{"from", ProjectTable},
				{"let", bson.D{{"pid", "$project_id"}}},
				{"pipeline", bson.A{
					bson.D {
						{
							"$match", bson.D{
								{"$expr", bson.D{
									{"$eq", bson.A{"$_id", "$$pid"}},
								}},
							},
						},
					},
				}},
				{"as", "project"},
			},
			},
		},
		//第二个连表是查对应的队列名称以及请求的接口配置
		{
			{
				"$lookup", bson.D{
					{"from", eventTypeConfigTable},
					{"let", bson.D{{"eventId", "$_id"}}},
					{"pipeline", bson.A {
						bson.D {
							{
								"$match", bson.D{
									{"$expr", bson.D{
										{"$eq", bson.A{"$event_id", "$$eventId"}},
									}},
								},
							},
						},
						bson.D {
							{
								"$lookup", bson.D {
									{"from", CallbackRequestConfigTable},
									{"let", bson.D{{"ecid", "$_id"}}},
									{
										"pipeline", bson.A {
											bson.D {
												{
													"$match", bson.D{
														{
															"$expr", bson.D{
																{"$eq", bson.A{"$parent_id", "$$ecid"}},
															},
														},
													},
												},
											},
										},
									},
									{"as", "crc"},
								},
							},
						},
					}},
					{"as", "ec"},
				},
			},
		},
		//unwind
		{
			{
				"$unwind", bson.D{
					{"path", "$project"},
					{"preserveNullAndEmptyArrays", false},
				},
			},
		},
		{
			{
				"$unwind", bson.D{
					{"path", "$ec"},
					{"preserveNullAndEmptyArrays", false},
				},
			},
		},
		{
			{
				"$unwind", bson.D{
					{"path", "$ec.crc"},
					{"preserveNullAndEmptyArrays", false},
				},
			},
		},
		{
			{
				"$project", project,
			},
		},
	}
	return
}

//根据事件消息队列类型获取对应的数据out字段
func getOutFields(eventType string) (project bson.D) {
	switch eventType {
		case "Single" :
			project = bson.D {
				//{"_id", 0},
				{"name" , 1},
				{"project.name" , 1},
				{"ec._id", 1},
				{"ec.queue_name", 1},
				{"ec.crc.callback_host", 1},
				{"ec.crc.callback_path", 1},
				{"ec.crc.callback_request_type", 1},
				{"ec.crc.callback_request_is_json", 1},
			}
		case "WorkQueues" :
			project = bson.D {
				//{"_id", 0},
				{"name" , 1},
				{"project.name" , 1},
				{"ec._id", 1},
				{"ec.queue_name", 1},
				{"ec.crc.callback_host", 1},
				{"ec.crc.callback_path", 1},
				{"ec.crc.callback_request_type", 1},
				{"ec.crc.callback_request_is_json", 1},
			}
		case "PublishSubscribe" :
			project = bson.D {
				//{"_id", 0},
				{"name" , 1},
				{"exchange_type" , 1},
				{"project.name" , 1},
				{"ec._id", 1},
				{"ec.queue_name", 1},
				{"ec.exchange_name", 1},
				{"ec.routing_key", 1},
				{"ec.headers", 1},
				{"ec.x_match", 1},
				{"ec.crc.callback_host", 1},
				{"ec.crc.callback_path", 1},
				{"ec.crc.callback_request_type", 1},
				{"ec.crc.callback_request_is_json", 1},
			}
	}
	return
}

//根据事件消息队列类型获取对应的配置数据表名
func getEventTypeConfigTable(eventType string) (tableName string) {
	switch eventType {
		case "Single" :
			tableName = SingleConfigTable
		case "WorkQueues":
			tableName = WorkQueuesConfigTable
		case "PublishSubscribe" :
			tableName = PublishSubscribeTable
	}
	return
}

//根据队列类型获取所有相对应的配置
func GetAllConfigsByType(typeStr string, eventId primitive.ObjectID) (results []bson.M) {
	collection := GetProjectEventCollection()
	var pipeline mongo.Pipeline
	pipeline = getPipeline(typeStr, eventId)
	cursor, _ := collection.Aggregate(context.TODO(), pipeline)
	err := cursor.All(context.TODO(), &results)
	if err != nil {}
	return
}

//根据id获取event数据
func GetEventsByIds(ids []primitive.ObjectID) ([]ProjectEvent, error) {
	var results []ProjectEvent
	var err error
	collection := GetProjectEventCollection()
	filter := bson.M{"_id": bson.M{"$in": ids}}
	res, decodeErr := collection.Find(context.TODO(), filter)
	if decodeErr != nil {
		err = decodeErr
	}
	if res != nil {
		err = res.All(context.TODO(), &results)
	}
	return results, err
}