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

//workQueues消息队列消费者协程
func WorkQueuesConsume(eventId primitive.ObjectID) {
	res := model.GetAllConfigsByType("WorkQueues", eventId)
	if len(res) > 0 {
		for _, item := range res {
			go workQueuesGoroutine(item)
		}
	}
}

func workQueuesGoroutine(item bson.M) {
	eventId := item["_id"].(primitive.ObjectID)
	var rabbitmqWorkQueues base.RabbitmqWorkQueues
	_, err := rabbitmqWorkQueues.InitConsoleWorkQueuesParams(item)
	if err != nil {
		common.FailOnError(err, "Failed to init params")
	}
	messages, err := rabbitmqWorkQueues.BasicConsume(item["ec"].(bson.M)["_id"].(primitive.ObjectID).String())
	common.FailOnError(err, "Failed to register a consumer")
	//chanConsumeMessages(messages, item)
	for {
		select {
		case chanEventId := <-config.CloseChan:
			if chanEventId == eventId {
				_ = rabbitmqWorkQueues.Channel.QueueUnbind(rabbitmqWorkQueues.GeneratedQueue, rabbitmqWorkQueues.GeneratedRoutingKey, rabbitmqWorkQueues.GeneratedQueue, amqp.Table{})
				_, _ = rabbitmqWorkQueues.Channel.QueueDelete(rabbitmqWorkQueues.GeneratedQueue, false, false, false)
				runtime.Goexit()
			} else {
				config.CloseChan <- chanEventId
				time.Sleep(500 * time.Millisecond)
			}
		case msg := <- messages:
			if msg.Acknowledger == nil {
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