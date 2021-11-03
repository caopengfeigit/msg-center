package console

import (
	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"msgCenter/base"
	"msgCenter/common"
	"msgCenter/config"
	"msgCenter/model"
	"runtime"
	"time"
)

//single消息队列消费者协程
func SingleConsume(eventId primitive.ObjectID) {
	res := model.GetAllConfigsByType("Single", eventId)
	if len(res) > 0 {
		for _, item := range res {
			go singleGoroutine(item)
		}
	}
}

//开协程起single类型队列的消费者
func singleGoroutine(item bson.M) {
	eventId := item["_id"].(primitive.ObjectID)
	var rabbitmqSingle base.RabbitmqSingle
	_, err := rabbitmqSingle.InitConsoleSingleParams(item)
	if err != nil {
		common.FailOnError(err, "Failed to init params")
	}
	messages, err := rabbitmqSingle.BasicConsume(item["ec"].(bson.M)["_id"].(primitive.ObjectID).String())
	common.FailOnError(err, "Failed to register a consumer")
	for {
		select {
			case chanEventId := <-config.CloseChan:
				if chanEventId == eventId {
					_ = rabbitmqSingle.Channel.QueueUnbind(rabbitmqSingle.GeneratedQueue, rabbitmqSingle.GeneratedRoutingKey, rabbitmqSingle.GeneratedQueue, amqp.Table{})
					_, _ = rabbitmqSingle.Channel.QueueDelete(rabbitmqSingle.GeneratedQueue, false, false, false)
					runtime.Goexit()
				} else {
					config.CloseChan <- chanEventId
					time.Sleep(500 * time.Millisecond)
				}
			case msg := <- messages:
				//获取配置数据
				configData := common.GetCallbackConfigByItem(item)
				common.TryAndRetryConsume(configData, string(msg.Body))
				//对消息进行应答，ack 中 multiple 参数是表示是否是多条消息进行应答
				if err := msg.Ack(false); err != nil {
					err := msg.Ack(false)
					log.Println(err)
				}
		}
	}
}

//消费消息
func chanConsumeMessages(messages <-chan amqp.Delivery, item bson.M) {
	/*forever := make(chan bool)
	go func() {
		for msg := range messages {
			//获取配置数据
			configData := common.GetCallbackConfigByItem(item)
			common.TryAndRetryConsume(configData, string(msg.Body))
			//对消息进行应答，ack 中 multiple 参数是表示是否是多条消息进行应答
			if err := msg.Ack(false); err != nil {
				err := msg.Ack(false)
				log.Println(err)
			}
		}
	}()
	<-forever*/
	/*for msg := range messages {
		//获取配置数据
		configData := common.GetCallbackConfigByItem(item)
		common.TryAndRetryConsume(configData, string(msg.Body))
		//对消息进行应答，ack 中 multiple 参数是表示是否是多条消息进行应答
		if err := msg.Ack(false); err != nil {
			err := msg.Ack(false)
			log.Println(err)
		}
	}*/
}