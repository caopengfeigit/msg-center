package model

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"msgCenter/config"
)

//配置列表
type ConfigList struct {
	No int
	EventId primitive.ObjectID
	ProjectName string
	EventName string
	EventType string
	ExchangeType string
}

/**
 * 获取配置列表数据及数据总条数
 */
func GetConfigList(projectName string, eventName string, page int) (list []ConfigList, total int) {
	//match条件
	var projectCond bson.E
	var eventCond bson.D
	//搜索业务线
	if projectName != "" {
		projectCond = bson.E{"name", projectName}
	}
	//搜索事件
	if eventName != "" {
		eventCond = bson.D{
			{
				"$match" , bson.D {
					{"name", eventName},
				},
			},
		}
	}
	
	//offset
	offset := (page - 1) * config.PageSize
	//获取查询列表数据pipeline
	pipeline := getConfigListPipeline(projectCond, eventCond, offset, config.PageSize, false)
	//获取查询列表数据总条数pipeline
	countPipeline := getConfigListPipeline(projectCond, eventCond, offset, config.PageSize, true)
	collection := GetProjectEventCollection()
	cursor, _ := collection.Aggregate(context.TODO(), pipeline)
	countCursor, _ := collection.Aggregate(context.TODO(), countPipeline)
	var results, countRes []bson.M
	err := cursor.All(context.TODO(), &results)
	countErr := countCursor.All(context.TODO(), &countRes)
	if err != nil || countErr != nil {}
	if len(countRes) > 0 {
		total = int(countRes[0]["total"].(int32))
	}
	
	//数据结构调整
	if len(results) > 0 {
		for key, val := range results {
			var item ConfigList
			item.No = (page - 1) * config.PageSize + key + 1
			item.ProjectName = val["project"].(bson.M)["name"].(string)
			item.EventId = val["_id"].(primitive.ObjectID)
			item.EventName = val["name"].(string)
			item.EventType = val["type"].(string)
			if val["exchange_type"] != nil {
				item.ExchangeType = val["exchange_type"].(string)
			} else {
				item.ExchangeType = "-"
			}
			list = append(list, item)
		}
	}
	return
}

/**
 * 获取配置列表pipeline条件
 */
func getConfigListPipeline(projectCond bson.E, eventCond bson.D, skip int, pageSize int, getCount bool) (pipeline mongo.Pipeline) {
	var countCond bson.D
	if getCount {
		countCond = bson.D {
			bson.E{"$count", "total"},
		}
	}
	if len(eventCond) > 0 {
		//event name match
		pipeline = append(pipeline, eventCond)
	}
	//project排序
	eventSort := bson.D{
		{
			"$sort" , bson.D {
				{"_id", -1},
			},
		},
	}
	//连表查询project数据
	joinProjectCond := bson.D {
		{
			"$lookup", bson.D {
				{"from", ProjectTable},
				{"let", bson.D{{"pid", "$project_id"}}},
				{"pipeline", bson.A{
					bson.D {
						{
							"$match", bson.D{
								{"$expr", bson.D{
									{"$eq", bson.A{"$_id", "$$pid"}},
								}},
								projectCond,
							},
						},
					},
				}},
				{"as", "project"},
			},
		},
	}
	//match过滤project为空的数据
	/*filterMatch := bson.D {
		{
			"$match", bson.D{
				{"project", bson.A{"$ne", bson.D{}}},
			},
		},
	}*/
	//unwind
	unwind := bson.D {
		{
			"$unwind", bson.D{
				{"path", "$project"},
				{"preserveNullAndEmptyArrays", false},
			},
		},
	}
	//project
	project := bson.D {
		{
			"$project", bson.D {
				{"_id", 1},
				{"name" , 1},
				{"type" , 1},
				{"exchange_type" , 1},
				{"project.name" , 1},
			},
		},
	}
	//分页limit
	limitCond := bson.D {
		{"$limit", pageSize},
	}
	//分页skip
	skipCond := bson.D {
		{"$skip", skip},
	}
	pipeline = append(pipeline, eventSort)
	pipeline = append(pipeline, joinProjectCond)
	pipeline = append(pipeline, unwind)
	pipeline = append(pipeline, project)
	
	//获取总条数，一定要放在skip和limit之间
	if len(countCond) > 0 {
		pipeline = append(pipeline, countCond)
	} else {
		//一定要先skip，再limit
		pipeline = append(pipeline, skipCond)
		pipeline = append(pipeline, limitCond)
	}
	
	return
}