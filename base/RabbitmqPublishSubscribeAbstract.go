package base

import (
	"github.com/streadway/amqp"
	"strconv"
)

type RabbitmqPublishSubscribeAbstract struct {
	RabbitmqInterface
	PublishSubscribeConfig
}

//声明交换器
func (rabbitmqPublishSubscribe *RabbitmqPublishSubscribeAbstract) ExchangeDeclare() (code int, err error) {
	//延时队列需要额外的参数
	args :=make(amqp.Table)
	if rabbitmqPublishSubscribe.IsDelay {
		args["x-delayed-type"] = amqp.ExchangeDirect
	}
	
	err = rabbitmqPublishSubscribe.Channel.ExchangeDeclare(
		rabbitmqPublishSubscribe.GeneratedExchange,
		rabbitmqPublishSubscribe.ExchangeType,
		true,          //持久化
		true,      //自动删除
		false,        //是否是内置交互器,(只能通过交换器将消息路由到此交互器，不能通过客户端发送消息
		false,
		args,
	)
	if err != nil {
		code = 20004
	}
	return
}

//声明队列
func (rabbitmqPublishSubscribe *RabbitmqPublishSubscribeAbstract) QueueDeclare(args amqp.Table) (code int, err error) {
	rabbitmqPublishSubscribe.DeclaredQueue, err = rabbitmqPublishSubscribe.Channel.QueueDeclare(
		rabbitmqPublishSubscribe.GeneratedQueue, // name
		true,   // durable
		true,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		args,     // arguments
	)
	if err != nil {
		code = 20002
	}
	return
}

//绑定队列
func (rabbitmqPublishSubscribe *RabbitmqPublishSubscribeAbstract) QueueBind() (code int, err error) {
	args :=make(amqp.Table)
	if rabbitmqPublishSubscribe.ExchangeType == amqp.ExchangeHeaders {
		args["x-match"] = rabbitmqPublishSubscribe.XMatch
		for key, value := range rabbitmqPublishSubscribe.Headers {
			args[key] = value
		}
	}
	err = rabbitmqPublishSubscribe.Channel.QueueBind(
		rabbitmqPublishSubscribe.GeneratedQueue, // queue name
		rabbitmqPublishSubscribe.GeneratedRoutingKey,     // routing key
		rabbitmqPublishSubscribe.GeneratedExchange, // exchange
		false,
		args)
	if err != nil {
		code = 20005
	}
	return
}

//发布一条消息
func (rabbitmqPublishSubscribe *RabbitmqPublishSubscribeAbstract) BasicPublish() (code int) {
	//publishing参数
	var arg amqp.Publishing
	arg.ContentType = "text/plain"
	arg.Body = []byte(rabbitmqPublishSubscribe.Message)
	
	//如果是延时队列，需要延时时间
	if rabbitmqPublishSubscribe.IsDelay {
		//判断参数是否有延迟时间
		if rabbitmqPublishSubscribe.DelayTime <= 0 {
			code = 20001
			return
		}
		arg.Headers = amqp.Table{
			"x-delay" : strconv.Itoa(rabbitmqPublishSubscribe.DelayTime * 1000),
		}
	}
	//如果exchange类型是header，需要加入headers参数
	if rabbitmqPublishSubscribe.ExchangeType == amqp.ExchangeHeaders {
		headers := make(amqp.Table)
		for key, value := range rabbitmqPublishSubscribe.Headers {
			headers[key] = value
		}
		arg.Headers = headers
	}
	
	err := rabbitmqPublishSubscribe.Channel.Publish(
		rabbitmqPublishSubscribe.GeneratedExchange,     // exchange
		rabbitmqPublishSubscribe.GeneratedRoutingKey, // routing key
		false,  // mandatory
		false,  // immediate
		arg,
	)
	if err != nil {
		code = 20003
	}
	return
}

//消费消息，⚠️注意：需要手动ack
func (rabbitmqPublishSubscribe *RabbitmqPublishSubscribeAbstract) BasicConsume(consumer string) (messages <-chan amqp.Delivery, err error) {
	messages, err = rabbitmqPublishSubscribe.Channel.Consume(
		rabbitmqPublishSubscribe.GeneratedQueue, // queue
		consumer,     // consumer
		false,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	return
}