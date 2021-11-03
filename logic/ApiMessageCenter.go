package logic

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"math"
	"msgCenter/config"
	"msgCenter/console"
	"msgCenter/db"
	"msgCenter/model"
	"strings"
	"time"
)

//获取配置列表数据
func GetList(projectName string, eventName string, page int) (list []model.ConfigList, pagination config.Pagination) {
	if page < 1 {
		page = 1
	}
	list, total := model.GetConfigList(projectName, eventName, page)
	pagination.Page = page
	pagination.PageSize = config.PageSize
	pagination.Total = total
	totalPage := float64(total) / float64(config.PageSize)
	pagination.TotalPage = int(math.Ceil(totalPage))
	return
}

//获取业务线列表数据
func GetProjectList(name string, businessLine string, page int) (list []model.ProjectList, pagination config.Pagination) {
	if page < 1 {
		page = 1
	}
	list, total := model.GetProjectList(name, businessLine, page)
	pagination.Page = page
	pagination.PageSize = config.PageSize
	pagination.Total = total
	totalPage := float64(total) / float64(config.PageSize)
	pagination.TotalPage = int(math.Ceil(totalPage))
	return
}

//删除事件及连带的数据
func DelEventConfig(id string) (success bool, code int) {
	idArr := strings.Split(id, ",")
	var ids []primitive.ObjectID
	if len(idArr) > 0 {
		for _, val := range idArr {
			id, _ := primitive.ObjectIDFromHex(val)
			ids = append(ids, id)
		}
	}
	if len(ids) > 0 {
		//先找到对应的event数据
		events, err := model.GetEventsByIds(ids)
		if err != nil {
			code = 30002
			return
		}
		
		//循环查询配置和回调配置数据进行删除（之所以循环查询，是因为不同的event对应的数据表不同）
		if len(events) > 0 {
			for _, event := range events {
				var collection *mongo.Collection
				//查找对应的config数据
				filter := bson.M{"event_id": event.Id}
				configCount := 0
				switch event.Type {
					case "Single" :
						collection = model.GetSingleConfigCollection()
						var configs []model.SingleConfig
						res, decodeErr := collection.Find(context.TODO(), filter)
						if decodeErr != nil {
							continue
						}
						if res != nil {
							err := res.All(context.TODO(), &configs)
							if err != nil {
								continue
							}
						}
						
						//循环single config查找回调配置数据
						var configIds, callbackIds []primitive.ObjectID
						if len(configs) > 0 {
							configCount = len(configs)
							for _, item := range configs {
								callbackRequests, code := model.GetCallbackRequestConfigsByParentId(item.Id)
								if code > 0 {
									continue
								}
								configIds = append(configIds, item.Id)
								if len(callbackRequests) > 0 {
									for _, request := range callbackRequests {
										callbackIds = append(callbackIds, request.Id)
									}
								}
							}
							
							//删除event、config、callbackRequest数据
							delConfigs(event.Id, configIds, callbackIds, collection)
						}
					case "WorkQueues":
						collection = model.GetWorkQueuesConfigCollection()
						var configs []model.WorkQueuesConfig
						res, decodeErr := collection.Find(context.TODO(), filter)
						if decodeErr != nil {
							continue
						}
						if res != nil {
							err := res.All(context.TODO(), &configs)
							if err != nil {
								continue
							}
						}
						
						//循环config查找回调配置数据
						var configIds, callbackIds []primitive.ObjectID
						if len(configs) > 0 {
							configCount = len(configs)
							for _, item := range configs {
								callbackRequests, code := model.GetCallbackRequestConfigsByParentId(item.Id)
								if code > 0 {
									continue
								}
								configIds = append(configIds, item.Id)
								if len(callbackRequests) > 0 {
									for _, request := range callbackRequests {
										callbackIds = append(callbackIds, request.Id)
									}
								}
							}
							
							//删除event、config、callbackRequest数据
							delConfigs(event.Id, configIds, callbackIds, collection)
						}
					case "PublishSubscribe" :
						collection = model.GetPublishSubscribeConfigCollection()
						var configs []model.PublishSubscribeConfig
						res, decodeErr := collection.Find(context.TODO(), filter)
						if decodeErr != nil {
							continue
						}
						if res != nil {
							err := res.All(context.TODO(), &configs)
							if err != nil {
								continue
							}
						}
						
						//循环config查找回调配置数据
						var configIds, callbackIds []primitive.ObjectID
						if len(configs) > 0 {
							configCount = len(configs)
							for _, item := range configs {
								callbackRequests, code := model.GetCallbackRequestConfigsByParentId(item.Id)
								if code > 0 {
									continue
								}
								configIds = append(configIds, item.Id)
								if len(callbackRequests) > 0 {
									for _, request := range callbackRequests {
										callbackIds = append(callbackIds, request.Id)
									}
								}
							}
							
							//删除event、config、callbackRequest数据
							delConfigs(event.Id, configIds, callbackIds, collection)
						}
				}
				//通过管道关闭对应的协程
				for i := 0; i< configCount; i++ {
					config.CloseChan <- event.Id
				}
			}
		}
	}
	success = true
	return
}

//事务删除对应的配置数据
func delConfigs(eventId primitive.ObjectID, configIds []primitive.ObjectID, callbackIds []primitive.ObjectID, configCollection *mongo.Collection) {
	ctx := context.Background()
	db.MongodbConn.UseSession(ctx, func(sessionContext mongo.SessionContext) error {
		err := sessionContext.StartTransaction()
		if err != nil {
			return err
		}
		
		//先删除event数据
		eventCollection := model.GetProjectEventCollection()
		_, delErr := eventCollection.DeleteOne(context.TODO(), bson.M{"_id": eventId})
		if delErr != nil {
			sessionContext.AbortTransaction(sessionContext)
			return delErr
		}
		
		//批量删除config数据
		_, delErr = configCollection.DeleteMany(context.TODO(), bson.M{"_id": bson.M{"$in": configIds}})
		if delErr != nil {
			sessionContext.AbortTransaction(sessionContext)
			return delErr
		}
		
		//批量删除回调request数据
		callbackCollection := model.GetCallBackRequestConfigCollection()
		_, delErr = callbackCollection.DeleteMany(context.TODO(), bson.M{"_id": bson.M{"$in": callbackIds}})
		if delErr != nil {
			sessionContext.AbortTransaction(sessionContext)
			return delErr
		}
		sessionContext.CommitTransaction(sessionContext)
		
		return nil
	})
}
//删除业务线及连带的数据
func DelProject(id string) (success bool, code int) {
	idArr := strings.Split(id, ",")
	var ids []primitive.ObjectID
	if len(idArr) > 0 {
		for _, val := range idArr {
			id, _ := primitive.ObjectIDFromHex(val)
			ids = append(ids, id)
		}
	}
	if len(ids) > 0 {
		//循环进行删除
		for _, id := range ids {
			//查找对应的event ids
			var events []model.ProjectEvent
			eventCollection := model.GetProjectEventCollection()
			cursor, err := eventCollection.Find(context.TODO(), bson.M{"project_id": id})
			if err != nil {
				continue
			}
			if cursor != nil {
				err := cursor.All(context.TODO(), &events)
				if err != nil {
					continue
				}
			}
			
			//整理出events ids
			var eventIds []string
			if len(events) > 0 {
				for _, event := range events {
					eventIds = append(eventIds, event.Id.Hex())
				}
			}
			
			//先删除项目
			projectCollection := model.GetProjectCollection()
			_, delErr := projectCollection.DeleteOne(context.TODO(), bson.M{"_id": id})
			if delErr != nil {
				continue
			} else {
				//删除对应的events以及其下的配置
				var eventIdsStr string
				if len(eventIds) > 0 {
					eventIdsStr = strings.Join(eventIds, ",")
					DelEventConfig(eventIdsStr)
				}
			}
		}
	}
	success = true
	return
}

//添加业务线
func AddProject(name string, businessLine string) (success bool, code int) {
	//先查对应的name有没有存在
	collection := model.GetProjectCollection()
	var project model.Project
	_ = collection.FindOne(context.TODO(), bson.M{"name": name}).Decode(&project)
	if project == (model.Project{}) {
		_, err := collection.InsertOne(context.TODO(), bson.M{"name": name, "business_line": businessLine})
		if err == nil {
			success = true
		} else {
			code = 40002
		}
	} else {
		code = 40001
	}
	return
}

//修改业务线
func EditProject(objectId primitive.ObjectID, name string, businessLine string) (success bool, code int) {
	//先查查有没有该业务线
	collection := model.GetProjectCollection()
	var project model.Project
	decodeErr := collection.FindOne(context.TODO(), bson.M{"_id": objectId}).Decode(&project)
	if decodeErr != nil || project == (model.Project{}) {
		code = 40003
	} else {
		//判断当前修改的name是否和原name一致，不一致需要查表是否已经存在该标识对应的业务线
		if name != project.Name {
			var findProject model.Project
			_ = collection.FindOne(context.TODO(), bson.M{"name": name}).Decode(&findProject)
			if findProject != (model.Project{}) {
				code = 40001
			}
		}
		//更新
		if name != project.Name || businessLine != project.BusinessLine {
			t := true
			_, err := collection.UpdateOne(
				context.TODO(),
				bson.M{"_id": objectId},
				bson.M{"$set": bson.M{"name": name, "business_line": businessLine}},
				&options.UpdateOptions{Upsert: &t},
			)
			if err != nil {
				code = 40004
			} else {
				success = true
				//消费者重启
				if name != project.Name {
					ReloadConsumersWhenEditProject(objectId)
				}
			}
		} else {
			success = true
		}
	}
	return
}

//获取配置详情
func GetConfigDetail(objectId primitive.ObjectID) (configDetail interface{}, code int) {
	//先根据eventId获取对应的event数据
	var event model.ProjectEvent
	eventCollection := model.GetProjectEventCollection()
	_ = eventCollection.FindOne(context.TODO(), bson.M{"_id": objectId}).Decode(&event)
	if event == (model.ProjectEvent{}) {
		code = 40005
	} else {
		//根据不同的event type 返回不同的配置详情数据
		switch event.Type {
			case "Single":
				configDetail = model.GetSingleDetail(objectId)
			case "WorkQueues":
				configDetail = model.GetWorkQueuesDetail(objectId)
			case "PublishSubscribe":
				configDetail = model.GetPublishSubscribeDetail(objectId, event.ExchangeType)
		}
	}
	return
}

//获取回调请求记录
func GetCallbackLogs(projectName string, eventKey string, startDate time.Time, endDate time.Time, page int) (list []model.CallbackRequestLog, pagination config.Pagination) {
	if page < 1 {
		page = 1
	}
	list, total := model.GetCallbackLogs(projectName, eventKey, startDate, endDate, page)
	pagination.Page = page
	pagination.PageSize = config.PageSize
	pagination.Total = total
	totalPage := float64(total) / float64(config.PageSize)
	pagination.TotalPage = int(math.Ceil(totalPage))
	return
}

//搜索业务线
func SearchProject(projectName string) (project model.Project) {
	project, _ = model.GetOneProjectByName(projectName)
	return
}

//搜索事件
func SearchEvent(projectId primitive.ObjectID, eventName string) (event model.ProjectEvent) {
	event, _ = model.GetOneProjectEventByIdName(projectId, eventName)
	return
}

//添加配置
func AddConfig(params config.AddConfigRequestParams) (success bool, code int) {
	//先判断有没有该业务线
	projectId, err := primitive.ObjectIDFromHex(params.ProjectId)
	if err != nil {
		code = 10001
		return
	}
	project := model.GetProjectById(projectId)
	if project == (model.Project{}) {
		code = 40003
		return
	}
	
	//判断传的事件有没有存在
	event, _ := model.GetOneProjectEventByIdName(projectId, params.EventName)
	if event != (model.ProjectEvent{}) {
		code = 40006
		return
	}
	
	//判断小队列类型是否合法
	if params.EventType != "Single" && params.EventType != "WorkQueues" && params.EventType != "PublishSubscribe" {
		code = 10001
	}
	
	//根据消息队列类型走不同的逻辑
	switch params.EventType {
		case "Single":
			valid := checkSingleParams(params)
			if !valid {
				code = 10001
				return
			}
			//数据入库
			res, eventId := addSingleConfig(projectId, params)
			if !res {
				code = 40002
				return
			}
			success = true
			//起消费者
			go console.SingleConsume(eventId)
		case "WorkQueues":
			valid := checkWorkQueuesParams(params)
			if !valid {
				code = 10001
				return
			}
			//数据入库
			res, eventId := addWrokQueuesConfig(projectId, params)
			if !res {
				code = 40002
				return
			}
			success = true
			//起消费者
			go console.WorkQueuesConsume(eventId)
		case "PublishSubscribe":
			valid := checkPublishSubscribeParams(params)
			if !valid {
				code = 10001
				return
			}
			//数据入库
			res, eventId := addPublishSubscribeConfig(projectId, params)
			if !res {
				code = 40002
				return
			}
			success = true
			go console.PublishSubscribeConsume(eventId)
	}
	
	return
}

//single参数校验
func checkSingleParams(params config.AddConfigRequestParams) bool {
	valid := true
	if params.QueueName == "" || params.CallbackHost == "" || params.CallbackPath == "" || (params.CallbackRequestType != "GET" && params.CallbackRequestType != "POST") {
		valid = false
	}
	return valid
}

//workqueues参数校验
func checkWorkQueuesParams(params config.AddConfigRequestParams) bool {
	valid := true
	if params.QueueName == "" || params.CallbackHost == "" || params.CallbackPath == "" || (params.CallbackRequestType != "GET" && params.CallbackRequestType != "POST") || params.WorkQueuesNum < 2 || params.WorkQueuesNum > 5 {
		valid = false
	}
	return valid
}

//publishSubscribe参数校验
func checkPublishSubscribeParams(params config.AddConfigRequestParams) bool {
	valid := true
	if (params.ExchangeType != "direct" && params.ExchangeType != "fanout" && params.ExchangeType != "headers" && params.ExchangeType != "topic" && params.ExchangeType != "x-delayed-message") || params.ExchangeName == "" {
		valid = false
	} else {
		switch params.ExchangeType {
			case "direct":
				if params.RoutingKey == "" || params.QueueName == "" || params.CallbackHost == "" || params.CallbackPath == "" || (params.CallbackRequestType != "GET" && params.CallbackRequestType != "POST") {
					valid = false
				}
				break
			case "fanout":
				queueNames := make(map[string]string)
				for _, item := range params.FanoutConfigs {
					if item.QueueName == "" || item.CallbackHost == "" || item.CallbackPath == "" || (item.CallbackRequestType != "GET" && item.CallbackRequestType != "POST") {
						valid = false
						break
					}
					if _, ok := queueNames[item.QueueName]; ok {
						valid = false
						break
					}
					queueNames[item.QueueName] = item.QueueName
				}
				break
			case "headers":
				if params.QueueName == "" || (params.XMatch != "any" && params.XMatch != "all") || params.CallbackHost == "" || params.CallbackPath == "" || (params.CallbackRequestType != "GET" && params.CallbackRequestType != "POST") {
					valid = false
				} else {
					headerKeys := make(map[string]string)
					for _, item := range params.Headers {
						if item.Key == "" || item.Value == "" {
							valid = false
							break
						}
						if _, ok := headerKeys[item.Key]; ok {
							valid = false
							break
						}
						headerKeys[item.Key] = item.Key
					}
				}
				break
			case "topic":
				routingKeys := make(map[string]string)
				for _, item := range params.Topic {
					if item.RoutingKey == "" || item.QueueName == "" || item.CallbackHost == "" || item.CallbackPath == "" || (params.CallbackRequestType != "GET" && params.CallbackRequestType != "POST") {
						valid = false
						break
					}
					if _, ok := routingKeys[item.RoutingKey]; ok {
						valid = false
						break
					}
					routingKeys[item.RoutingKey] = item.RoutingKey
				}
				break
			case "x-delayed-message":
				if params.RoutingKey == "" || params.QueueName == "" || params.CallbackHost == "" || params.CallbackPath == "" || (params.CallbackRequestType != "GET" && params.CallbackRequestType != "POST") {
					valid = false
				}
				break
		}
	}
	return valid
}

//single数据入库
func addSingleConfig(projectId primitive.ObjectID, params config.AddConfigRequestParams) (res bool, eventId primitive.ObjectID) {
	res = true
	ctx := context.Background()
	err := db.MongodbConn.UseSession(ctx, func(sessionContext mongo.SessionContext) error {
		err := sessionContext.StartTransaction()
		if err != nil {
			return err
		}
		
		//先添加event事件数据
		eventCollection := model.GetProjectEventCollection()
		insertEventRes, err := eventCollection.InsertOne(context.TODO(), bson.M{"project_id": projectId, "name": params.EventName, "type": params.EventType})
		if err != nil {
			_ = sessionContext.AbortTransaction(sessionContext)
			return err
		}
		eventId = insertEventRes.InsertedID.(primitive.ObjectID)
		
		//再添加single_config数据
		singleConfigCollection := model.GetSingleConfigCollection()
		insertSingleConfigRes, err := singleConfigCollection.InsertOne(context.TODO(), bson.M{"event_id": eventId, "queue_name": params.QueueName})
		if err != nil {
			_ = sessionContext.AbortTransaction(sessionContext)
			return err
		}
		singleConfigId := insertSingleConfigRes.InsertedID.(primitive.ObjectID)
		
		//最后添加callback config 数据
		callbackCollection := model.GetCallBackRequestConfigCollection()
		_, err = callbackCollection.InsertOne(context.TODO(), bson.M{
			"parent_id": singleConfigId,
			"callback_host": params.CallbackHost,
			"callback_path": params.CallbackPath,
			"callback_request_type": params.CallbackRequestType,
			"callback_request_is_json": params.CallbackRequestIsJson,
		})
		if err != nil {
			_ = sessionContext.AbortTransaction(sessionContext)
			return err
		}
		
		_ = sessionContext.CommitTransaction(sessionContext)
		
		return nil
	})
	if err != nil {
		res = false
	}
	return
}

//workqueues数据入库
func addWrokQueuesConfig(projectId primitive.ObjectID, params config.AddConfigRequestParams) (res bool, eventId primitive.ObjectID) {
	res = true
	ctx := context.Background()
	err := db.MongodbConn.UseSession(ctx, func(sessionContext mongo.SessionContext) error {
		err := sessionContext.StartTransaction()
		if err != nil {
			return err
		}
		
		//先添加event事件数据
		eventCollection := model.GetProjectEventCollection()
		insertEventRes, err := eventCollection.InsertOne(context.TODO(), bson.M{"project_id": projectId, "name": params.EventName, "type": params.EventType})
		if err != nil {
			_ = sessionContext.AbortTransaction(sessionContext)
			return err
		}
		eventId = insertEventRes.InsertedID.(primitive.ObjectID)
		
		//再添加work_queues_config数据
		workQueuesConfigCollection := model.GetWorkQueuesConfigCollection()
		for i := 0; i < params.WorkQueuesNum; i++ {
			insertWorkQueuesConfigRes, err := workQueuesConfigCollection.InsertOne(context.TODO(), bson.M{"event_id": eventId, "queue_name": params.QueueName})
			if err != nil {
				_ = sessionContext.AbortTransaction(sessionContext)
				return err
			}
			workQueuesConfigId := insertWorkQueuesConfigRes.InsertedID.(primitive.ObjectID)
			
			//最后添加callback config 数据
			callbackCollection := model.GetCallBackRequestConfigCollection()
			_, err = callbackCollection.InsertOne(context.TODO(), bson.M{
				"parent_id": workQueuesConfigId,
				"callback_host": params.CallbackHost,
				"callback_path": params.CallbackPath,
				"callback_request_type": params.CallbackRequestType,
				"callback_request_is_json": params.CallbackRequestIsJson,
			})
			if err != nil {
				_ = sessionContext.AbortTransaction(sessionContext)
				return err
			}
		}
		
		_ = sessionContext.CommitTransaction(sessionContext)
		
		return nil
	})
	if err != nil {
		res = false
	}
	return
}

//publishsubscribe数据入库
func addPublishSubscribeConfig(projectId primitive.ObjectID, params config.AddConfigRequestParams) (res bool, eventId primitive.ObjectID) {
	res = true
	ctx := context.Background()
	err := db.MongodbConn.UseSession(ctx, func(sessionContext mongo.SessionContext) error {
		err := sessionContext.StartTransaction()
		if err != nil {
			return err
		}
		
		//先添加event事件数据
		eventCollection := model.GetProjectEventCollection()
		insertEventRes, err := eventCollection.InsertOne(context.TODO(), bson.M{"project_id": projectId, "name": params.EventName, "type": params.EventType, "exchange_type": params.ExchangeType})
		if err != nil {
			_ = sessionContext.AbortTransaction(sessionContext)
			return err
		}
		eventId = insertEventRes.InsertedID.(primitive.ObjectID)
		
		//根据exchange type向不同的collection添加数据
		publishSubscribeCollection := model.GetPublishSubscribeConfigCollection()
		switch params.ExchangeType {
			case "direct":
				fallthrough
			case "x-delayed-message":
				insertPublishSubscribeConfigRes, err := publishSubscribeCollection.InsertOne(context.TODO(), bson.M{
					"event_id": eventId,
					"exchange_name": params.ExchangeName,
					"routing_key": params.RoutingKey,
					"queue_name": params.QueueName,
				})
				if err != nil {
					_ = sessionContext.AbortTransaction(sessionContext)
					return err
				}
				publishSubscribeConfigId := insertPublishSubscribeConfigRes.InsertedID.(primitive.ObjectID)
				
				//最后添加callback config 数据
				callbackCollection := model.GetCallBackRequestConfigCollection()
				_, err = callbackCollection.InsertOne(context.TODO(), bson.M{
					"parent_id": publishSubscribeConfigId,
					"callback_host": params.CallbackHost,
					"callback_path": params.CallbackPath,
					"callback_request_type": params.CallbackRequestType,
					"callback_request_is_json": params.CallbackRequestIsJson,
				})
				if err != nil {
					_ = sessionContext.AbortTransaction(sessionContext)
					return err
				}
			case "fanout":
				for _, item := range params.FanoutConfigs {
					insertPublishSubscribeConfigRes, err := publishSubscribeCollection.InsertOne(context.TODO(), bson.M{
						"event_id": eventId,
						"exchange_name": params.ExchangeName,
						"queue_name": item.QueueName,
					})
					if err != nil {
						_ = sessionContext.AbortTransaction(sessionContext)
						return err
					}
					publishSubscribeConfigId := insertPublishSubscribeConfigRes.InsertedID.(primitive.ObjectID)
					
					//最后添加callback config 数据
					callbackCollection := model.GetCallBackRequestConfigCollection()
					_, err = callbackCollection.InsertOne(context.TODO(), bson.M{
						"parent_id": publishSubscribeConfigId,
						"callback_host": item.CallbackHost,
						"callback_path": item.CallbackPath,
						"callback_request_type": item.CallbackRequestType,
						"callback_request_is_json": item.CallbackRequestIsJson,
					})
					if err != nil {
						_ = sessionContext.AbortTransaction(sessionContext)
						return err
					}
				}
			case "headers":
				//格式化headers
				formatHeaders := make(bson.M)
				for _, item := range params.Headers {
					formatHeaders[item.Key] = item.Value
				}
				insertPublishSubscribeConfigRes, err := publishSubscribeCollection.InsertOne(context.TODO(), bson.M{
					"event_id": eventId,
					"exchange_name": params.ExchangeName,
					"queue_name": params.QueueName,
					"headers": formatHeaders,
					"x_match": params.XMatch,
				})
				if err != nil {
					_ = sessionContext.AbortTransaction(sessionContext)
					return err
				}
				publishSubscribeConfigId := insertPublishSubscribeConfigRes.InsertedID.(primitive.ObjectID)
				
				//最后添加callback config 数据
				callbackCollection := model.GetCallBackRequestConfigCollection()
				_, err = callbackCollection.InsertOne(context.TODO(), bson.M{
					"parent_id": publishSubscribeConfigId,
					"callback_host": params.CallbackHost,
					"callback_path": params.CallbackPath,
					"callback_request_type": params.CallbackRequestType,
					"callback_request_is_json": params.CallbackRequestIsJson,
				})
				if err != nil {
					_ = sessionContext.AbortTransaction(sessionContext)
					return err
				}
			case "topic":
				for _, item := range params.Topic {
					insertPublishSubscribeConfigRes, err := publishSubscribeCollection.InsertOne(context.TODO(), bson.M{
						"event_id": eventId,
						"exchange_name": params.ExchangeName,
						"routing_key": item.RoutingKey,
						"queue_name": item.QueueName,
					})
					if err != nil {
						_ = sessionContext.AbortTransaction(sessionContext)
						return err
					}
					publishSubscribeConfigId := insertPublishSubscribeConfigRes.InsertedID.(primitive.ObjectID)
					
					//最后添加callback config 数据
					callbackCollection := model.GetCallBackRequestConfigCollection()
					_, err = callbackCollection.InsertOne(context.TODO(), bson.M{
						"parent_id": publishSubscribeConfigId,
						"callback_host": item.CallbackHost,
						"callback_path": item.CallbackPath,
						"callback_request_type": item.CallbackRequestType,
						"callback_request_is_json": item.CallbackRequestIsJson,
					})
					if err != nil {
						_ = sessionContext.AbortTransaction(sessionContext)
						return err
					}
				}
		}
		
		_ = sessionContext.CommitTransaction(sessionContext)
		
		return nil
	})
	if err != nil {
		res = false
	}
	return
}

//修改业务线name重启消费者
func ReloadConsumersWhenEditProject(objectId primitive.ObjectID) {
	//先找出event列表
	var events []model.ProjectEvent
	collection := model.GetProjectEventCollection()
	filter := bson.D{{"project_id", objectId}}
	res, _ := collection.Find(context.TODO(), filter)
	if res != nil {
		_ = res.All(context.TODO(), &events)
	}
	if len(events) > 0 {
		for _, event := range events {
			var count int
			var res int64
			var err error
			var configCollection *mongo.Collection
			//先根据消息队列类型，查有多少条配置数据
			switch event.Type {
			case "Single":
				configCollection = model.GetSingleConfigCollection()
			case "WorkQueues":
				configCollection = model.GetWorkQueuesConfigCollection()
			case "PublishSubscribe":
				configCollection = model.GetPublishSubscribeConfigCollection()
			}
			res, err = configCollection.CountDocuments(context.TODO(), bson.M{"event_id": event.Id})
			if err == nil {
				count = int(res)
			}
			//先通知管道进行关闭修改之前的消费者
			for i := 0; i < count; i++ {
				config.CloseChan <- event.Id
			}
			time.Sleep(500 * time.Millisecond)
			//再重新启动协程
			switch event.Type {
			case "Single":
				if !event.Id.IsZero() {
					go console.SingleConsume(event.Id)
				}
			case "WorkQueues":
				if !event.Id.IsZero() {
					go console.WorkQueuesConsume(event.Id)
				}
			case "PublishSubscribe":
				if !event.Id.IsZero() {
					go console.PublishSubscribeConsume(event.Id)
				}
			}
		}
	}
}