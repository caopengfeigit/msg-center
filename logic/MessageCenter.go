package logic

import (
	"github.com/streadway/amqp"
	"msgCenter/base"
	"msgCenter/common"
	"msgCenter/config"
	"msgCenter/model"
)

func HandleMessage(requestParams config.RequestParams) (code int) {
	//配置分割
	configKey := common.SplitParams(requestParams.ConfigKey)
	if len(configKey) != 2 {
		code = 10001
		return
	}
	
	//从数据库中查找是否有对应的项目
	var project model.Project
	project, code = model.GetOneProjectByName(configKey[0])
	if code > 0 {
		return
	}
	
	//查找对应的项目下是否有对应的配置
	var projectEvent model.ProjectEvent
	projectEvent, code = model.GetOneProjectEventByIdName(project.Id, configKey[1])
	if code > 0 {
		return
	}
	
	//如果是发布订阅模式且交换机类型是topic时，必须传入routingKey
	if projectEvent.ExchangeType == amqp.ExchangeTopic {
		if requestParams.RoutingKey == "" {
			code = 20001
			return
		}
	}
	
	//按照不同的消息类型去找不同的struct
	switch projectEvent.Type {
		case "Single" :
			var rabbitmqSingle base.RabbitmqSingle
			//参数初始化
			code = rabbitmqSingle.InitSingleParams(requestParams, project, projectEvent)
			if code > 0 {
				return
			}
			//发布消息
			code = rabbitmqSingle.BasicPublish()
			return
		case "WorkQueues" :
			var rabbitmqWorkQueues base.RabbitmqWorkQueues
			//参数初始化
			code = rabbitmqWorkQueues.InitWorkQueuesParams(requestParams, project, projectEvent)
			if code > 0 {
				return
			}
			//发布消息
			code = rabbitmqWorkQueues.BasicPublish()
			return
		case "PublishSubscribe" :
			var rabbitmqPublishSubscribe base.RabbitmqPublishSubscribe
			//参数初始化
			code = rabbitmqPublishSubscribe.InitPublishSubscribeParams(requestParams, project, projectEvent)
			if code > 0 {
				return
			}
			//发布消息
			code = rabbitmqPublishSubscribe.BasicPublish()
			return
		default :
			code = 30002
			return
	}
	
	return
}
