package config

type errorMessage map[int]string

var ErrorMessage = errorMessage {
	//参数验证类 1开头
	10001 : "parameter error",
	
	//程序处理类 2开头
	20001 : "program error",
	20002 : "failed to declare a queue",
	20003 : "failed to publish message",
	20004 : "failed to declare an exchange",
	20005 : "failed to bind a queue",
	
	//数据库类 3开头
	30001 : "save error",
	30002 : "data does not exists",
	30003 : "data already exists",
	30004 : "database error",
	
	//业务报错信息类 4开头
	40001 : "已存在该标识对应的业务线",//business line of this name has already exist
	40002 : "添加失败",//add failed
	40003 : "不存在该业务线",//this business line doesn't exist
	40004 : "编辑失败",//edit failed
	40005 : "不存在该事件",//this event doesn't exist
	40006 : "当前业务线下已存在该事件",//event has already exist
}