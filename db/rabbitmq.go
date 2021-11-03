package db

import (
	"github.com/streadway/amqp"
	"msgCenter/config"
)

var (
	RabbitmqConn *amqp.Connection
	Channel *amqp.Channel
	mqErr, channelErr error
)

func RabbitmqInit() {
	//可连接多个服务
	defaultMqConnect()
}

func defaultMqConnect() {
	RabbitmqConn, mqErr = amqp.Dial("amqp://" + config.RabbitmqInc.User + ":" + config.RabbitmqInc.Password + "@" + config.RabbitmqInc.Host + ":" + config.RabbitmqInc.Port)
	if mqErr != nil{
		panic(mqErr.Error())
	}
	//defer RabbitmqConn.Close()
	
	//开启一个channel
	Channel, channelErr = RabbitmqConn.Channel()
	if channelErr != nil{
		panic(channelErr.Error())
	}
	//defer Channel.Close()
}

//关闭rabbitmq资源
func CloseRabbitmqResource() {
	if RabbitmqConn != nil && !RabbitmqConn.IsClosed() {
		defer RabbitmqConn.Close()
		defer Channel.Close()
	}
}