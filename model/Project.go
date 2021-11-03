package model

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"msgCenter/config"
	"msgCenter/db"
)

var ProjectTable = "project"

type Project struct {
	Id primitive.ObjectID `bson:"_id"`
	Name string `bson:"name"`
	BusinessLine string `bson:"business_line"`
}

//获取mongo生成的collection
func GetProjectCollection() *mongo.Collection {
	return db.MongodbConn.Database(config.MongodbInc.MongoDB).Collection(ProjectTable)
}

//通过name获取一个项目
func GetOneProjectByName(name string) (project Project, code int) {
	collection := db.MongodbConn.Database(config.MongodbInc.MongoDB).Collection(ProjectTable)
	filter := bson.D{{"name", name}}
	decodeErr := collection.FindOne(context.TODO(), filter).Decode(&project)
	if decodeErr != nil {
		if project == (Project{}) {
			code = 30002
		} else {
			code = 30004
		}
		return
	}
	return
}

//通过id获取项目
func GetProjectById(id primitive.ObjectID) (project Project) {
	collection := db.MongodbConn.Database(config.MongodbInc.MongoDB).Collection(ProjectTable)
	filter := bson.D{{"_id", id}}
	_ = collection.FindOne(context.TODO(), filter).Decode(&project)
	return
}