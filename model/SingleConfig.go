package model

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"msgCenter/config"
	"msgCenter/db"
)

var SingleConfigTable = "single_config"

type SingleConfig struct {
	Id primitive.ObjectID `bson:"_id"`
	EventId primitive.ObjectID `bson:"event_id"`
	QueueName string `bson:"queue_name"`
}

//获取mongo生成的collection
func GetSingleConfigCollection() *mongo.Collection {
	return db.MongodbConn.Database(config.MongodbInc.MongoDB).Collection(SingleConfigTable)
}

//通过eventId获取一条single配置
func GetSingleConfig(eventId primitive.ObjectID) (singleConfig SingleConfig, code int) {
	collection := GetSingleConfigCollection()
	filter := bson.D{{"event_id", eventId}}
	decodeErr := collection.FindOne(context.TODO(), filter).Decode(&singleConfig)
	if decodeErr != nil {
		if singleConfig == (SingleConfig{}) {
			code = 30002
		} else {
			code = 30004
		}
	}
	return
}