package model

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

//single配置详情
type SingleDetail struct {
	QueueName string `bson:"queue_name"`
	CallbackHost string `bson:"callback_host"`
	CallbackPath string `bson:"callback_path"`
	CallbackRequestType string `bson:"callback_request_type"`
	CallbackRequestIsJson bool `bson:"callback_request_is_json"`
}

//workqueues配置详情
type WorkQueuesDetail struct {
	QueueName string `bson:"queue_name"`
	CallbackHost string `bson:"callback_host"`
	CallbackPath string `bson:"callback_path"`
	CallbackRequestType string `bson:"callback_request_type"`
	CallbackRequestIsJson bool `bson:"callback_request_is_json"`
}

//publishsubscribe配置详情
type PublishSubscribeDetail struct {
	QueueName string `bson:"queue_name"`
	ExchangeName string `bson:"exchange_name"`
	RoutingKey string `bson:"routing_key"`
	Headers bson.M `bson:"headers"`
	XMatch string `bson:"x_match"`
	CallbackHost string `bson:"callback_host"`
	CallbackPath string `bson:"callback_path"`
	CallbackRequestType string `bson:"callback_request_type"`
	CallbackRequestIsJson bool `bson:"callback_request_is_json"`
}

//获取single配置详情
func GetSingleDetail(eventId primitive.ObjectID) (singleDetail SingleDetail) {
	//获取single config
	var singleConfig SingleConfig
	configCollection := GetSingleConfigCollection()
	_ = configCollection.FindOne(context.TODO(), bson.M{"event_id": eventId}).Decode(&singleConfig)
	if singleConfig != (SingleConfig{}) {
		//先赋值队列名
		singleDetail.QueueName = singleConfig.QueueName
		//再找callback config
		var callbackConfig CallbackRequestConfig
		callbackCollection := GetCallBackRequestConfigCollection()
		_ = callbackCollection.FindOne(context.TODO(), bson.M{"parent_id": singleConfig.Id}).Decode(&callbackConfig)
		if callbackConfig != (CallbackRequestConfig{}) {
			//赋值其它配置
			singleDetail.CallbackHost = callbackConfig.CallbackHost
			singleDetail.CallbackPath = callbackConfig.CallbackPath
			singleDetail.CallbackRequestType = callbackConfig.CallbackRequestType
			singleDetail.CallbackRequestIsJson = callbackConfig.CallbackRequestIsJson
		}
	}
	return
}

//获取workqueues配置详情
func GetWorkQueuesDetail(eventId primitive.ObjectID) (workQueuesDetails []WorkQueuesDetail) {
	workQueuesCollection := GetWorkQueuesConfigCollection()
	pipeline := getConfigDetailPipeline(eventId, "WorkQueues")
	cursor, _ := workQueuesCollection.Aggregate(context.TODO(), pipeline)
	var details []bson.M
	_ = cursor.All(context.TODO(), &details)
	if len(details) > 0 {
		var detail WorkQueuesDetail
		for _, item := range details {
			detail.QueueName = item["queue_name"].(string)
			detail.CallbackHost = item["callback"].(bson.M)["callback_host"].(string)
			detail.CallbackPath = item["callback"].(bson.M)["callback_path"].(string)
			detail.CallbackRequestType = item["callback"].(bson.M)["callback_request_type"].(string)
			detail.CallbackRequestIsJson = item["callback"].(bson.M)["callback_request_is_json"].(bool)
			workQueuesDetails = append(workQueuesDetails, detail)
		}
	}
	return
}

//获取publishsubscribe配置详情
func GetPublishSubscribeDetail(eventId primitive.ObjectID, exchangeType string) (publishSubscribeDetails []PublishSubscribeDetail) {
	publishSubscribeCollection := GetPublishSubscribeConfigCollection()
	pipeline := getConfigDetailPipeline(eventId, "PublishSubscribe")
	cursor, _ := publishSubscribeCollection.Aggregate(context.TODO(), pipeline)
	var details []bson.M
	_ = cursor.All(context.TODO(), &details)
	if len(details) > 0 {
		var detail PublishSubscribeDetail
		for _, item := range details {
			detail.QueueName = item["queue_name"].(string)
			detail.ExchangeName = item["exchange_name"].(string)
			switch exchangeType {
				case "direct":
					detail.RoutingKey = item["routing_key"].(string)
				case "x-delayed-message":
					detail.RoutingKey = item["routing_key"].(string)
				case "topic":
					detail.RoutingKey = item["routing_key"].(string)
				case "headers":
					detail.Headers = item["headers"].(bson.M)
					detail.XMatch = item["x_match"].(string)
			}
			detail.CallbackHost = item["callback"].(bson.M)["callback_host"].(string)
			detail.CallbackPath = item["callback"].(bson.M)["callback_path"].(string)
			detail.CallbackRequestType = item["callback"].(bson.M)["callback_request_type"].(string)
			detail.CallbackRequestIsJson = item["callback"].(bson.M)["callback_request_is_json"].(bool)
			publishSubscribeDetails = append(publishSubscribeDetails, detail)
		}
	}
	return
}

//获取查询配置详情pipeline
func getConfigDetailPipeline(eventId primitive.ObjectID, eventType string) mongo.Pipeline {
	fileds := getConfigDetailProject(eventType)
	pipeline := mongo.Pipeline {
		//match条件
		{
			{
				"$match" , bson.D {
					{"event_id", eventId},
				},
			},
		},
		//连callback表查询数据
		{
			{
				"$lookup", bson.D{
					{"from", CallbackRequestConfigTable},
					{"let", bson.D{{"parentId", "$_id"}}},
					{"pipeline", bson.A{
						bson.D {
							{
								"$match", bson.D{
									{"$expr", bson.D{
										{"$eq", bson.A{"$parent_id", "$$parentId"}},
									}},
								},
							},
						},
					}},
					{"as", "callback"},
				},
			},
		},
		//unwind
		{
			{
				"$unwind", bson.D{
					{"path", "$callback"},
					{"preserveNullAndEmptyArrays", false},
				},
			},
		},
		//project
		{
			{
				"$project", fileds,
			},
		},
	}
	return pipeline
}

//获取配置详情展示的字段
func getConfigDetailProject(eventType string) (project bson.D) {
	switch eventType {
		case "WorkQueues" :
			project = bson.D {
				{"queue_name", 1},
				{"callback.callback_host", 1},
				{"callback.callback_path", 1},
				{"callback.callback_request_type", 1},
				{"callback.callback_request_is_json", 1},
			}
		case "PublishSubscribe" :
			project = bson.D {
				{"queue_name", 1},
				{"exchange_name", 1},
				{"routing_key", 1},
				{"headers", 1},
				{"x_match", 1},
				{"callback.callback_host", 1},
				{"callback.callback_path", 1},
				{"callback.callback_request_type", 1},
				{"callback.callback_request_is_json", 1},
			}
	}
	return
}