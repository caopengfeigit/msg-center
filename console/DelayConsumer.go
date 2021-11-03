package console

import (
	"github.com/Braveheart7854/rabbitmqPool"
	"github.com/streadway/amqp"
	"msgCenter/common"
	"msgCenter/config"
)

//延时队列消息消费协程
func DelayConsume() {
	queueName := config.RabbitmqSettingConfig["delayQueueName"].(string)
	exchange := config.RabbitmqSettingConfig["delayExchangeName"].(string)
	channel, _ := rabbitmqPool.AmqpServer.GetChannel()
	
	//声明交换机
	exchangeErr := channel.ExchangeDeclare(
		exchange,
		amqp.ExchangeDirect,   //交换机模式为direct
		true,          //持久化
		false,      //自动删除
		true,        //是否是内置交互器,(只能通过交换器将消息路由到此交互器，不能通过客户端发送消息
		false,
		nil,
	)
	common.FailOnError(exchangeErr, "Failed to declare a exchange")
	
	//声明队列
	queue, queueErr := channel.QueueDeclare(
		queueName, // name
		true,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	common.FailOnError(queueErr, "Failed to declare a queue")
	
	//将交换器和队列绑定
	bindErr := channel.QueueBind(
		queue.Name,
		"",
		exchange,
		false,
		nil,
	)
	common.FailOnError(bindErr, "Failed to bind a queue to an exchange")
	
	messages, err := channel.Consume(
		queue.Name, // queue
		"",     // consumer
		false,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	common.FailOnError(err, "Failed to register a consumer")
	
	forever := make(chan bool)
	go func() {
		for msg := range messages {
			//从头信息中获取配置数据
			if msg.Headers["projectName"] != nil && msg.Headers["eventKey"] != nil || msg.Headers["requestHost"] != nil || msg.Headers["requestPath"] != nil || msg.Headers["requestType"] != nil {
				item := make(map[string]interface{})
				item["projectName"] = msg.Headers["projectName"].(string)
				item["eventKey"] = msg.Headers["eventKey"].(string)
				item["queueName"] = msg.Headers["queueName"].(string)
				item["requestHost"] = msg.Headers["requestHost"].(string)
				item["requestPath"] = msg.Headers["requestPath"].(string)
				item["requestType"] = msg.Headers["requestType"].(string)
				item["requestIsJson"] = msg.Headers["requestIsJson"].(bool)
				
				common.TryAndRetryConsume(item, string(msg.Body))
			}
			//对消息进行应答，ack 中 multiple 参数是表示是否是多条消息进行应答
			if err := msg.Ack(false); err != nil {
				msg.Ack(false)
			}
		}
	}()
	<-forever
}
