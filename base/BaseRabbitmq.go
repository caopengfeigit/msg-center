package base

import (
	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/bson"
)

//rabbitmq公共配置结构
type BaseRabbitMqConfig struct {
	//channel池拿到的channel
	Channel *amqp.Channel
	BaseProjectConfig
	BaseProjectRabbimtConfig
}

//项目对应的rabbitmq配置
type BaseProjectRabbimtConfig struct {
	//生成routingkey或queue时的前缀
	Prefix string
	//收到的消息
	Message string
	//交换器名
	Exchange string
	//针对每个配置生成一个对应的exchange
	GeneratedExchange string
	//routingKey
	RoutingKey string
	//针对每个配置生成一个对应的routingkey
	GeneratedRoutingKey string
	//队列名
	QueueName string
	//headers模式下的headers参数
	Headers bson.M
	//headers模式下x-match的匹配规则
	XMatch string
	//针对每个配置生成一个对应的队列queue名（字符串）
	GeneratedQueue string
	//通过rabbitmq生成的队列queue
	DeclaredQueue amqp.Queue
}

//项目配置结构
type BaseProjectConfig struct {
	//消息队列的类型
	MessageType string
	//项目名
	ProjectName string
	//事件key
	EventKey string
}

//延时队列配置
type BaseDelayConfig struct {
	//是否是延时队列
	IsDelay bool
	//延时时长，单位是秒
	DelayTime int
}

//callback配置结构
type BaseCallbackConfig struct {
	//回调host
	CallBackHost string
	//回调接口路径
	CallbackPath string
	//回调请求方式
	CallbackRequestType string
	//是否以json格式发送
	CallbackRequestIsJson bool
}

//queues组配置 qname=>{回调配置}
type BaseQueuesConfig map[string]struct {
	BaseCallbackConfig
}

//routingkey组配置 routingkey=>{队列和回调配置}
type BaseRoutingKeyConfig map[string]struct {
	BaseCallbackConfig
}

//single类型配置
type SingleConfig struct {
	BaseRabbitMqConfig
	BaseCallbackConfig
}

//workerQueues类型配置
type WorkerQueuesConfig struct {
	//沿用single类型的配置
	SingleConfig
}

//publish-subscribe类型配置
type PublishSubscribeConfig struct {
	//交换器类型
	ExchangeType string
	BaseRabbitMqConfig
	BaseDelayConfig
	//queues相关配置
	BaseQueuesConfig
	BaseCallbackConfig
}
