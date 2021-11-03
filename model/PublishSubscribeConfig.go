package model

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"msgCenter/config"
	"msgCenter/db"
)

var PublishSubscribeTable = "publish_subscribe_config"

type PublishSubscribeConfig struct {
	Id primitive.ObjectID `bson:"_id"`
	EventId primitive.ObjectID `bson:"event_id"`
	ExchangeName string `bson:"exchange_name"`
	RoutingKey string `bson:"routing_key"`
	QueueName string `bson:"queue_name"`
	Headers interface{} `bson:"headers"`
	XMatch string `bson:"x_match"`
}

//获取mongo生成的collection
func GetPublishSubscribeConfigCollection() *mongo.Collection {
	return db.MongodbConn.Database(config.MongodbInc.MongoDB).Collection(PublishSubscribeTable)
}

//通过eventId获取一条single配置
func GetPublishSubscribeConfig(eventId primitive.ObjectID) (publishSubscribeConfig PublishSubscribeConfig, code int) {
	collection := GetPublishSubscribeConfigCollection()
	filter := bson.D{{"event_id", eventId}}
	decodeErr := collection.FindOne(context.TODO(), filter).Decode(&publishSubscribeConfig)
	if decodeErr != nil {
		if publishSubscribeConfig == (PublishSubscribeConfig{}) {
			code = 30002
		} else {
			code = 30004
		}
	}
	return
}