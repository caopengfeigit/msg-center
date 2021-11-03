package db

import (
	"context"
	"fmt"
	
	"msgCenter/config"
	"strings"
	
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	MongodbConn *mongo.Client
	MongodbErr error
)

func MongodbInit() {
	//可连接多个服务
	defaultMongodbConnect()
}

//获取mongodb连接uri
func getMongodbConnUri() string {
	var builder strings.Builder
	builder.WriteString("mongodb://")
	builder.WriteString(config.MongodbInc.MongoUser)
	builder.WriteString(":")
	builder.WriteString(config.MongodbInc.MongoPwd)
	builder.WriteString("@")
	builder.WriteString(config.MongodbInc.MongoHost)
	builder.WriteString(":")
	builder.WriteString(config.MongodbInc.MongoPort)
	builder.WriteString("/")
	builder.WriteString(config.MongodbInc.MongoDB)
	return builder.String()
}

func defaultMongodbConnect() {
	fmt.Println("执行初始化mongodb服务连接")
	uri := getMongodbConnUri()
	fmt.Println(uri)
	clientOptions := options.Client().ApplyURI(uri)
	MongodbConn, MongodbErr = mongo.Connect(context.TODO(), clientOptions)
	if MongodbErr != nil {
		panic(MongodbErr.Error())
	}
	//defer MongodbConn.Disconnect(context.TODO())
}