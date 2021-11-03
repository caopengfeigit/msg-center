package base

import (
	"github.com/Braveheart7854/rabbitmqPool"
	"go.mongodb.org/mongo-driver/bson"
	"msgCenter/config"
	"msgCenter/db"
	"msgCenter/model"
)

type RabbitmqWorkQueues struct {
	RabbitmqWorkQueuesAbstract
}

//初始化参数
func (rabbitmqWorkQueues *RabbitmqWorkQueues) InitWorkQueuesParams(params config.RequestParams, project model.Project, projectEvent model.ProjectEvent) (code int) {
	db.RabbitmqInit()
	//数据库查找对应的配置并赋值
	var workQueuesConfig model.WorkQueuesConfig
	workQueuesConfig, code = model.GetWorkQueuesConfig(projectEvent.Id)
	if code > 0 {
		return
	}
	
	//数据库查找回调配置
	var callbackRequestConfig model.CallbackRequestConfig
	callbackRequestConfig, code = model.GetOneCallbackRequestConfigByParentId(workQueuesConfig.Id)
	if code > 0 {
		return
	}
	
	//参数赋值
	//rabbitmqSingle.Channel, _ = rabbitmqPool.AmqpServer.GetChannel()
	rabbitmqWorkQueues.Channel = db.Channel
	rabbitmqWorkQueues.MessageType = projectEvent.Type
	rabbitmqWorkQueues.ProjectName = project.Name
	rabbitmqWorkQueues.EventKey = projectEvent.Name
	rabbitmqWorkQueues.Prefix = project.Name + "_" + projectEvent.Name + "_"
	rabbitmqWorkQueues.Message = params.Message
	rabbitmqWorkQueues.QueueName = workQueuesConfig.QueueName
	rabbitmqWorkQueues.CallBackHost = callbackRequestConfig.CallbackHost
	rabbitmqWorkQueues.CallbackPath = callbackRequestConfig.CallbackPath
	rabbitmqWorkQueues.CallbackRequestType = callbackRequestConfig.CallbackRequestType
	rabbitmqWorkQueues.CallbackRequestIsJson = callbackRequestConfig.CallbackRequestIsJson
	rabbitmqWorkQueues.GeneratedQueue = rabbitmqWorkQueues.Prefix + rabbitmqWorkQueues.QueueName
	
	return
}

//消费者初始化参数
func (rabbitmqWorkQueues *RabbitmqWorkQueues) InitConsoleWorkQueuesParams(item bson.M) (code int, err error) {
	//参数赋值
	rabbitmqWorkQueues.Channel, _ = rabbitmqPool.AmqpServer.GetChannel()
	rabbitmqWorkQueues.MessageType = "Single"
	rabbitmqWorkQueues.ProjectName = item["project"].(bson.M)["name"].(string)
	rabbitmqWorkQueues.EventKey = item["name"].(string)
	rabbitmqWorkQueues.Prefix = rabbitmqWorkQueues.ProjectName + "_" + rabbitmqWorkQueues.EventKey + "_"
	rabbitmqWorkQueues.QueueName = item["ec"].(bson.M)["queue_name"].(string)
	rabbitmqWorkQueues.CallBackHost = item["ec"].(bson.M)["crc"].(bson.M)["callback_host"].(string)
	rabbitmqWorkQueues.CallbackPath = item["ec"].(bson.M)["crc"].(bson.M)["callback_path"].(string)
	rabbitmqWorkQueues.CallbackRequestType = item["ec"].(bson.M)["crc"].(bson.M)["callback_request_type"].(string)
	rabbitmqWorkQueues.CallbackRequestIsJson = item["ec"].(bson.M)["crc"].(bson.M)["callback_request_is_json"].(bool)
	rabbitmqWorkQueues.GeneratedQueue = rabbitmqWorkQueues.Prefix + rabbitmqWorkQueues.QueueName
	//声明队列
	code, err = rabbitmqWorkQueues.QueueDeclare(nil)
	return
}