package model

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"msgCenter/config"
)

//配置列表
type ProjectList struct {
	No int
	Id primitive.ObjectID `bson:"_id"`
	Name string `bson:"name"`
	BusinessLine string `bson:"business_line"`
}

/**
 * 获取配置列表数据及数据总条数
 */
func GetProjectList(name string, businessLine string, page int) (list []ProjectList, total int) {
	//查询条件
	filter := bson.M{}
	if len(name) > 0 {
		filter["name"] = name
	}
	if len(businessLine) > 0 {
		filter["business_line"] = bson.M{"$regex": primitive.Regex{Pattern: businessLine}}
	}
	
	//分页
	skip := (page - 1) * config.PageSize
	findOptions := &options.FindOptions{}
	findOptions.SetLimit(int64(config.PageSize))
	findOptions.SetSkip(int64(skip))
	findOptions.SetSort(bson.M{"_id": -1})
	
	//查询数据列表
	collection := GetProjectCollection()
	res, decodeErr := collection.Find(context.TODO(), filter, findOptions)
	if decodeErr == nil && res != nil {
		err := res.All(context.TODO(), &list)
		if err == nil {
			for key, item := range list {
				item.No = (page - 1) * config.PageSize + key + 1
				list[key] = item
			}
		}
	}
	
	//数据总数
	count, err := collection.CountDocuments(context.TODO(), filter)
	if err == nil {
		total = int(count)
	}
	
	return
}