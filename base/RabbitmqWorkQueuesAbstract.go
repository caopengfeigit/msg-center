package base

import "github.com/streadway/amqp"

type RabbitmqWorkQueuesAbstract struct {
	RabbitmqInterface
	WorkerQueuesConfig
}

//声明队列
func (rabbitmqWorkQueues *RabbitmqWorkQueuesAbstract) QueueDeclare(args amqp.Table) (code int, err error) {
	rabbitmqWorkQueues.DeclaredQueue, err = rabbitmqWorkQueues.Channel.QueueDeclare(
		rabbitmqWorkQueues.GeneratedQueue, // name
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
func (rabbitmqWorkQueues *RabbitmqWorkQueuesAbstract) BasicPublish() (code int) {
	//publishing参数
	var arg amqp.Publishing
	arg.ContentType = "text/plain"
	arg.Body = []byte(rabbitmqWorkQueues.Message)
	
	err := rabbitmqWorkQueues.Channel.Publish(
		"",     // exchange
		rabbitmqWorkQueues.GeneratedQueue, // routing key
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
func (rabbitmqWorkQueues *RabbitmqWorkQueuesAbstract) BasicConsume(consumer string) (messages <-chan amqp.Delivery, err error) {
	messages, err = rabbitmqWorkQueues.Channel.Consume(
		rabbitmqWorkQueues.GeneratedQueue, // queue
		consumer,     // consumer
		false,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	return
}