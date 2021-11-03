package base

import (
	"github.com/Braveheart7854/rabbitmqPool"
	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/bson"
	"msgCenter/config"
	"msgCenter/db"
	"msgCenter/model"
)

type RabbitmqPublishSubscribe struct {
	RabbitmqPublishSubscribeAbstract
}

//初始化参数
func (publishSubscribe *RabbitmqPublishSubscribe) InitPublishSubscribeParams(params config.RequestParams, project model.Project, projectEvent model.ProjectEvent) (code int) {
	db.RabbitmqInit()
	//数据库查找对应的配置并赋值
	var publishSubscribeConfig model.PublishSubscribeConfig
	publishSubscribeConfig, code = model.GetPublishSubscribeConfig(projectEvent.Id)
	if code > 0 {
		return
	}
	
	//数据库查找回调配置
	var callbackRequestConfig model.CallbackRequestConfig
	callbackRequestConfig, code = model.GetOneCallbackRequestConfigByParentId(publishSubscribeConfig.Id)
	if code > 0 {
		return
	}
	
	//参数赋值
	//rabbitmqSingle.Channel, _ = rabbitmqPool.AmqpServer.GetChannel()
	publishSubscribe.Channel = db.Channel
	publishSubscribe.MessageType = projectEvent.Type
	publishSubscribe.ProjectName = project.Name
	publishSubscribe.EventKey = projectEvent.Name
	publishSubscribe.Prefix = project.Name + "_" + projectEvent.Name + "_"
	//交换机类型如果是x-delayed-message则是延时队列
	if projectEvent.ExchangeType == config.RabbitmqSettingConfig["delayExchangeType"].(string) {
		publishSubscribe.IsDelay = true
	} else {
		publishSubscribe.IsDelay = false
	}
	publishSubscribe.Exchange = publishSubscribeConfig.ExchangeName
	publishSubscribe.ExchangeType = projectEvent.ExchangeType
	publishSubscribe.RoutingKey = publishSubscribeConfig.RoutingKey
	//交换机类型是topic的情况下，使用参数传过来的topic
	if publishSubscribe.ExchangeType == amqp.ExchangeTopic {
		publishSubscribe.RoutingKey = params.RoutingKey
	}
	publishSubscribe.QueueName = publishSubscribeConfig.QueueName
	publishSubscribe.Message = params.Message
	publishSubscribe.DelayTime = params.DelayTime
	publishSubscribe.CallBackHost = callbackRequestConfig.CallbackHost
	publishSubscribe.CallbackPath = callbackRequestConfig.CallbackPath
	publishSubscribe.CallbackRequestType = callbackRequestConfig.CallbackRequestType
	publishSubscribe.CallbackRequestIsJson = callbackRequestConfig.CallbackRequestIsJson
	publishSubscribe.GeneratedExchange = publishSubscribe.Prefix + publishSubscribeConfig.ExchangeName
	publishSubscribe.GeneratedQueue = publishSubscribe.Prefix + publishSubscribeConfig.QueueName
	if publishSubscribe.RoutingKey != "" {
		//topic模式下有通配符的存在，不加前缀
		if publishSubscribe.ExchangeType == amqp.ExchangeTopic {
			publishSubscribe.GeneratedRoutingKey = publishSubscribe.RoutingKey
		} else {
			publishSubscribe.GeneratedRoutingKey = publishSubscribe.Prefix + publishSubscribe.RoutingKey
		}
	}
	if publishSubscribe.ExchangeType == amqp.ExchangeHeaders {
		publishSubscribe.Headers = params.Headers
	}
	
	return
}

//消费者初始化参数
func (publishSubscribe *RabbitmqPublishSubscribe) InitConsolePublishSubscribeParams(item bson.M) (code int, err error) {
	//参数赋值
	publishSubscribe.Channel, _ = rabbitmqPool.AmqpServer.GetChannel()
	publishSubscribe.MessageType = "PublishSubscribe"
	publishSubscribe.ProjectName = item["project"].(bson.M)["name"].(string)
	publishSubscribe.EventKey = item["name"].(string)
	publishSubscribe.Prefix = publishSubscribe.ProjectName + "_" + publishSubscribe.EventKey + "_"
	//根据交换机的类型确定是否是延时队列
	if item["exchange_type"].(string) == config.RabbitmqSettingConfig["delayExchangeType"].(string) {
		publishSubscribe.IsDelay = true
	} else {
		publishSubscribe.IsDelay = false
	}
	publishSubscribe.Exchange = item["ec"].(bson.M)["exchange_name"].(string)
	publishSubscribe.ExchangeType = item["exchange_type"].(string)
	if item["ec"].(bson.M)["routing_key"] != nil {
		publishSubscribe.RoutingKey = item["ec"].(bson.M)["routing_key"].(string)
	}
	publishSubscribe.QueueName = item["ec"].(bson.M)["queue_name"].(string)
	publishSubscribe.CallBackHost = item["ec"].(bson.M)["crc"].(bson.M)["callback_host"].(string)
	publishSubscribe.CallbackPath = item["ec"].(bson.M)["crc"].(bson.M)["callback_path"].(string)
	publishSubscribe.CallbackRequestType = item["ec"].(bson.M)["crc"].(bson.M)["callback_request_type"].(string)
	publishSubscribe.CallbackRequestIsJson = item["ec"].(bson.M)["crc"].(bson.M)["callback_request_is_json"].(bool)
	publishSubscribe.GeneratedExchange = publishSubscribe.Prefix + publishSubscribe.Exchange
	publishSubscribe.GeneratedQueue = publishSubscribe.Prefix + publishSubscribe.QueueName
	if publishSubscribe.RoutingKey != "" {
		//topic模式下有通配符的存在，不加前缀
		if publishSubscribe.ExchangeType == amqp.ExchangeTopic {
			publishSubscribe.GeneratedRoutingKey = publishSubscribe.RoutingKey
		} else {
			publishSubscribe.GeneratedRoutingKey = publishSubscribe.Prefix + publishSubscribe.RoutingKey
		}
	}
	if publishSubscribe.ExchangeType == amqp.ExchangeHeaders {
		publishSubscribe.Headers = item["ec"].(bson.M)["headers"].(bson.M)
		publishSubscribe.XMatch = item["ec"].(bson.M)["x_match"].(string)
	}
	
	//声明延时交换机
	code, _ = publishSubscribe.ExchangeDeclare()
	if code > 0 {
		return
	}
	//声明队列
	code, _ = publishSubscribe.QueueDeclare(nil)
	if code > 0 {
		return
	}
	//绑定exchange和queue
	code, _ = publishSubscribe.QueueBind()
	
	return
}