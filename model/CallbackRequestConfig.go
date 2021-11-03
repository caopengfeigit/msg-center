package model

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"msgCenter/config"
	"msgCenter/db"
)

var CallbackRequestConfigTable = "callback_request_config"

type CallbackRequestConfig struct {
	Id primitive.ObjectID `bson:"_id"`
	ParentId string `bson:"parent_id"`
	CallbackHost string `bson:"callback_host"`
	CallbackPath string `bson:"callback_path"`
	CallbackRequestType string `bson:"callback_request_type"`
	CallbackRequestIsJson bool `bson:"callback_request_is_json"`
}

//获取mongo生成的collection
func GetCallBackRequestConfigCollection() *mongo.Collection {
	return db.MongodbConn.Database(config.MongodbInc.MongoDB).Collection(CallbackRequestConfigTable)
}

//通过parent_id获取一条回调配置
func GetOneCallbackRequestConfigByParentId(parentId primitive.ObjectID) (callbackRequestConfig CallbackRequestConfig, code int) {
	collection := GetCallBackRequestConfigCollection()
	filter := bson.D{{"parent_id", parentId}}
	decodeErr := collection.FindOne(context.TODO(), filter).Decode(&callbackRequestConfig)
	if decodeErr != nil {
		if callbackRequestConfig == (CallbackRequestConfig{}) {
			code = 30002
		} else {
			code = 30004
		}
	}
	return
}

//通过parent_id获取所有回调配置
func GetCallbackRequestConfigsByParentId(parentId primitive.ObjectID) (callbackRequestConfigs []CallbackRequestConfig, code int) {
	collection := GetCallBackRequestConfigCollection()
	filter := bson.D{{"parent_id", parentId}}
	res, decodeErr := collection.Find(context.TODO(), filter)
	if decodeErr != nil {
		if res == nil {
			code = 30002
		} else {
			code = 30004
		}
	}
	if res != nil {
		err := res.All(context.TODO(), &callbackRequestConfigs)
		if err != nil {
			if len(callbackRequestConfigs) == 0 {
				code = 30002
			} else {
				code = 30004
			}
		}
	}
	return
}