package main

import (
	"github.com/Braveheart7854/rabbitmqPool"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"msgCenter/console"
	"msgCenter/db"
	"msgCenter/router"
	"net/http"
)

func main() {
	//连接rabbitmq服务
	rabbitmqPool.InitAmqp()
	//连接mongodb服务
	db.MongodbInit()
	// 创建静态资源服务
	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/", fs)
	//设置gin
	gin.SetMode(gin.ReleaseMode) //设置debug模式 (release,debug.test)
	r := gin.Default()
	r.Delims("<<<", ">>>")
	r.LoadHTMLGlob("templates/**/*")
	r.StaticFS("/public", http.Dir("./public"))
	router.InitRouter(r)
	
	eventId, _ := primitive.ObjectIDFromHex("")
	//single消息队列启动
	go console.SingleConsume(eventId)
	//workQueues消息队列启动
	go console.WorkQueuesConsume(eventId)
	//publish subscribe消息队列启动
	go console.PublishSubscribeConsume(eventId)
	err := r.Run(":9501") // listen and serve on 0.0.0.0:9501
	if err != nil {
		panic(err)
	}
}
