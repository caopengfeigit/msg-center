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

func PublishSubscribeConsume(eventId primitive.ObjectID) {
	res := model.GetAllConfigsByType("PublishSubscribe", eventId)
	if len(res) > 0 {
		for _, item := range res {
			go publishSubscribeGoroutine(item)
		}
	}
}

//开协程起publish subscribe类型队列的消费者
func publishSubscribeGoroutine(item bson.M) {
	eventId := item["_id"].(primitive.ObjectID)
	var rabbitmqPublishSubscribe base.RabbitmqPublishSubscribe
	_, err := rabbitmqPublishSubscribe.InitConsolePublishSubscribeParams(item)
	if err != nil {
		common.FailOnError(err, "Failed to init params")
	}
	//消费
	messages, err := rabbitmqPublishSubscribe.BasicConsume(item["ec"].(bson.M)["_id"].(primitive.ObjectID).String())
	common.FailOnError(err, "Failed to register a consumer")
	//chanConsumeMessages(messages, item)
	for {
		select {
		case chanEventId := <-config.CloseChan:
			if chanEventId == eventId {
				_ = rabbitmqPublishSubscribe.Channel.QueueUnbind(rabbitmqPublishSubscribe.GeneratedQueue, rabbitmqPublishSubscribe.GeneratedRoutingKey, rabbitmqPublishSubscribe.GeneratedQueue, amqp.Table{})
				_, _ = rabbitmqPublishSubscribe.Channel.QueueDelete(rabbitmqPublishSubscribe.GeneratedQueue, false, false, false)
				runtime.Goexit()
			} else {
				config.CloseChan <- chanEventId
				time.Sleep(500 * time.Millisecond)
			}
		case msg := <- messages:
			if msg.Acknowledger == nil {
				_ = rabbitmqPublishSubscribe.Channel.QueueUnbind(rabbitmqPublishSubscribe.GeneratedQueue, rabbitmqPublishSubscribe.GeneratedRoutingKey, rabbitmqPublishSubscribe.GeneratedQueue, amqp.Table{})
				_, _ = rabbitmqPublishSubscribe.Channel.QueueDelete(rabbitmqPublishSubscribe.GeneratedQueue, false, false, false)
				runtime.Goexit()
			}
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
