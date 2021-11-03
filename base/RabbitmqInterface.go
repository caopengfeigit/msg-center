package base

//rabbitmq基础接口
type RabbitmqInterface interface {
	//声明一个队列
	QueueDeclare()
	
	//发布消息
	BasicPublish()
	
	//消费者消费消息
	BasicConsume()
}
