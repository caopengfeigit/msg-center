package config

import "go.mongodb.org/mongo-driver/bson"

//接受生产者生产的消息
type RequestParams struct {
	//这里因为兼容老系统，不进行require操作：binding:"required"
	Message string `form:"message"`
	ConfigKey string `form:"configKey"`
	RoutingKey string `form:"routingKey"`
	DelayTime int `form:"delayTime"`
	Headers bson.M `form:"headers"`
}

//添加配置参数
type AddConfigRequestParams struct {
	ProjectId string `json:"projectId" binding:"required"`
	EventName string `json:"eventName" binding:"required"`
	EventType string `json:"eventType" binding:"required"`
	WorkQueuesNum int `json:"workQueuesNum"`
	ExchangeType string `json:"exchangeType"`
	ExchangeName string `json:"exchangeName"`
	RoutingKey string `json:"routingKey"`
	QueueName string `json:"queueName"`
	CallbackHost string `json:"callbackHost"`
	CallbackPath string `json:"callbackPath"`
	CallbackRequestType string `json:"callbackRequestType"`
	CallbackRequestIsJson bool `json:"callbackRequestIsJson"`
	XMatch string `json:"xMatch"`
	FanoutConfigs []fanoutConfig `json:"fanoutConfigs"`
	Headers []headers `json:"headers"`
	Topic []topic `json:"topic"`
}

type fanoutConfig struct {
	QueueName string `json:"queueName"`
	CallbackHost string `json:"callbackHost"`
	CallbackPath string `json:"callbackPath"`
	CallbackRequestType string `json:"callbackRequestType"`
	CallbackRequestIsJson bool `json:"callbackRequestIsJson"`
}

type headers struct {
	Key string `json:"key"`
	Value string `json:"value"`
}

type topic struct {
	RoutingKey string `json:"routingKey"`
	QueueName string `json:"queueName"`
	CallbackHost string `json:"callbackHost"`
	CallbackPath string `json:"callbackPath"`
	CallbackRequestType string `json:"callbackRequestType"`
	CallbackRequestIsJson bool `json:"callbackRequestIsJson"`
}