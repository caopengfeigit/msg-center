package config

var RabbitmqSettingConfig = map[string]interface{} {
	//延时队列交换器类型
	"delayExchangeType" : "x-delayed-message",
	//延时队列交换器
	"delayExchangeName" : "delay_exchange",
	//延时队列名
	"delayQueueName" : "delay_queue",
	//消息消费失败时重试次数
	"tryReConsumeNum" : 2,
}