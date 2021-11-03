package base

import (
	"github.com/Braveheart7854/rabbitmqPool"
	"go.mongodb.org/mongo-driver/bson"
	"msgCenter/config"
	"msgCenter/db"
	"msgCenter/model"
)

type RabbitmqSingle struct {
	RabbitmqSingleAbstract
}

//初始化参数
func (rabbitmqSingle *RabbitmqSingle) InitSingleParams(params config.RequestParams, project model.Project, projectEvent model.ProjectEvent) (code int) {
	db.RabbitmqInit()
	//数据库查找对应的配置并赋值
	var singleConfig model.SingleConfig
	singleConfig, code = model.GetSingleConfig(projectEvent.Id)
	if code > 0 {
		return
	}
	
	//数据库查找回调配置
	var callbackRequestConfig model.CallbackRequestConfig
	callbackRequestConfig, code = model.GetOneCallbackRequestConfigByParentId(singleConfig.Id)
	if code > 0 {
		return
	}
	
	//参数赋值
	//rabbitmqSingle.Channel, _ = rabbitmqPool.AmqpServer.GetChannel()
	rabbitmqSingle.Channel = db.Channel
	rabbitmqSingle.MessageType = projectEvent.Type
	rabbitmqSingle.ProjectName = project.Name
	rabbitmqSingle.EventKey = projectEvent.Name
	rabbitmqSingle.Prefix = project.Name + "_" + projectEvent.Name + "_"
	rabbitmqSingle.Message = params.Message
	rabbitmqSingle.QueueName = singleConfig.QueueName
	rabbitmqSingle.CallBackHost = callbackRequestConfig.CallbackHost
	rabbitmqSingle.CallbackPath = callbackRequestConfig.CallbackPath
	rabbitmqSingle.CallbackRequestType = callbackRequestConfig.CallbackRequestType
	rabbitmqSingle.CallbackRequestIsJson = callbackRequestConfig.CallbackRequestIsJson
	rabbitmqSingle.GeneratedQueue = rabbitmqSingle.Prefix + rabbitmqSingle.QueueName
	
	return
}

//消费者初始化参数
func (rabbitmqSingle *RabbitmqSingle) InitConsoleSingleParams(item bson.M) (code int, err error) {
	//参数赋值
	rabbitmqSingle.Channel, _ = rabbitmqPool.AmqpServer.GetChannel()
	rabbitmqSingle.MessageType = "Single"
	rabbitmqSingle.ProjectName = item["project"].(bson.M)["name"].(string)
	rabbitmqSingle.EventKey = item["name"].(string)
	rabbitmqSingle.Prefix = rabbitmqSingle.ProjectName + "_" + rabbitmqSingle.EventKey + "_"
	rabbitmqSingle.QueueName = item["ec"].(bson.M)["queue_name"].(string)
	rabbitmqSingle.CallBackHost = item["ec"].(bson.M)["crc"].(bson.M)["callback_host"].(string)
	rabbitmqSingle.CallbackPath = item["ec"].(bson.M)["crc"].(bson.M)["callback_path"].(string)
	rabbitmqSingle.CallbackRequestType = item["ec"].(bson.M)["crc"].(bson.M)["callback_request_type"].(string)
	rabbitmqSingle.CallbackRequestIsJson = item["ec"].(bson.M)["crc"].(bson.M)["callback_request_is_json"].(bool)
	rabbitmqSingle.GeneratedQueue = rabbitmqSingle.Prefix + rabbitmqSingle.QueueName
	//声明队列
	code, err = rabbitmqSingle.QueueDeclare(nil)
	return
}
