package base

import (
	"github.com/streadway/amqp"
)

type RabbitmqSingleAbstract struct {
	RabbitmqInterface
	SingleConfig
}

//声明队列
func (rabbitmqSingle *RabbitmqSingleAbstract) QueueDeclare(args amqp.Table) (code int, err error) {
	rabbitmqSingle.DeclaredQueue, err = rabbitmqSingle.Channel.QueueDeclare(
		rabbitmqSingle.GeneratedQueue, // name
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

//发布一条消息
func (rabbitmqSingle *RabbitmqSingleAbstract) BasicPublish() (code int) {
	//publishing参数
	var arg amqp.Publishing
	arg.ContentType = "text/plain"
	arg.Body = []byte(rabbitmqSingle.Message)
	
	err := rabbitmqSingle.Channel.Publish(
		"",     // exchange
		rabbitmqSingle.GeneratedQueue, // routing key
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
func (rabbitmqSingle *RabbitmqSingleAbstract) BasicConsume(consumer string) (messages <-chan amqp.Delivery, err error) {
	messages, err = rabbitmqSingle.Channel.Consume(
		rabbitmqSingle.GeneratedQueue, // queue
		consumer,     // consumer
		false,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	return
}