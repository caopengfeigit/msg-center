package config

//rabbitmq配置
type RabbitmqConfig struct {
	Host string
	Port string //端口
	User string //用户名
	Password string //密码
}

//mongodb配置
type MongodbConfig struct {
	MongoHost string
	MongoPort string
	MongoDB string
	MongoUser string
	MongoPwd string
}

var RabbitmqInc = RabbitmqConfig {
	Host: "localhost",
	Port: "5672",
	User: "guest",
	Password: "guest",
}

var MongodbInc = MongodbConfig{
	MongoHost : "localhost",
	MongoPort : "27017",
	MongoDB : "msg_center",
	MongoUser : "msg_center",
	MongoPwd : "msg_center",
}
