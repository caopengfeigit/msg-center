package model

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"msgCenter/config"
	"msgCenter/db"
)

var WorkQueuesConfigTable = "work_queues_config"

type WorkQueuesConfig struct {
	Id primitive.ObjectID `bson:"_id"`
	EventId primitive.ObjectID `bson:"event_id"`
	QueueName string `bson:"queue_name"`
}

//获取mongo生成的collection
func GetWorkQueuesConfigCollection() *mongo.Collection {
	return db.MongodbConn.Database(config.MongodbInc.MongoDB).Collection(WorkQueuesConfigTable)
}

//通过eventId获取一条single配置
func GetWorkQueuesConfig(eventId primitive.ObjectID) (workQueuesConfig WorkQueuesConfig, code int) {
	collection := GetWorkQueuesConfigCollection()
	filter := bson.D{{"event_id", eventId}}
	decodeErr := collection.FindOne(context.TODO(), filter).Decode(&workQueuesConfig)
	if decodeErr != nil {
		if workQueuesConfig == (WorkQueuesConfig{}) {
			code = 30002
		} else {
			code = 30004
		}
	}
	return
}
