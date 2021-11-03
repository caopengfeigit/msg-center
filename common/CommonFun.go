package common

import (
	"encoding/json"
	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"msgCenter/config"
	"msgCenter/model"
	"strings"
	"time"
)

func FailOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

//分割参数
func SplitParams(params string) []string {
	return strings.Split(params, ".")
}

//消费者将消息发送给对应的接口消费方（msg从rabbitmq拿出来之后是byte类型，记得转成string之后再传过来）
func SendMessageToApi(item map[string]interface{}, msg string) bool {
	res := true
	var cstZone = time.FixedZone("CST", 8*3600)
	createdAt := time.Now().In(cstZone)
	//回调记录
	logData := &model.CallbackRequestLog{
		Id : primitive.NewObjectID(),
		ProjectName : item["projectName"].(string),
		EventKey : item["eventKey"].(string),
		QueueName : item["queueName"].(string),
		RequestHost : item["requestHost"].(string),
		RequestPath : item["requestPath"].(string),
		RequestType : item["requestType"].(string),
		CallbackRequestIsJson : item["requestIsJson"].(bool),
		RequestData : msg,
		RequestRes : true,
		RequestError : "",
		RequestStatus : 0,
		RequestResponse : "",
		CreatedAt : createdAt,
		CreatedAtStr : createdAt.Format("2006-01-02 15:04:05"),
	}
	
	url := item["requestHost"].(string) + item["requestPath"].(string)
	var status int
	var err error
	var responseBody string
	
	req := fasthttp.AcquireRequest()
	if item["requestIsJson"].(bool) {
		req.Header.SetContentType("application/json")
	}
	req.Header.SetMethod(item["requestType"].(string))
	req.SetRequestURI(url)
	data := map[string]string{"message" : msg}
	var requestBody []byte
	requestBody, err = json.Marshal(data)
	if err == nil {
		req.SetBody(requestBody)
		resp := fasthttp.AcquireResponse()
		err = fasthttp.Do(req, resp)
		status = resp.StatusCode()
		responseBody = string(resp.Body())
	}
	logData.RequestStatus = status
	logData.RequestResponse = responseBody
	if err != nil || status != fasthttp.StatusOK {
		res = false
		logData.RequestRes = false
		if err != nil {
			logData.RequestError = err.Error()
		}
	}
	//记录回调数据
	go model.InsertLog(logData)
	
	return res
}

//根据查询数据获取消息队列配置数据
func GetCallbackConfigByItem(item bson.M) map[string]interface{} {
	result := make(map[string]interface{})
	//项目名
	result["projectName"] = item["project"].(bson.M)["name"].(string)
	//事件名
	result["eventKey"] = item["name"].(string)
	//队列名
	result["queueName"] = item["ec"].(bson.M)["queue_name"].(string)
	//请求的接口host
	result["requestHost"] = item["ec"].(bson.M)["crc"].(bson.M)["callback_host"].(string)
	//请求的接口路径
	result["requestPath"] = item["ec"].(bson.M)["crc"].(bson.M)["callback_path"].(string)
	//请求类型
	result["requestType"] = item["ec"].(bson.M)["crc"].(bson.M)["callback_request_type"].(string)
	//post请求是否是json格式
	result["requestIsJson"] = item["ec"].(bson.M)["crc"].(bson.M)["callback_request_is_json"].(bool)
	return result
}

//根据配置数据进行回调和错误回调重试
func TryAndRetryConsume(configData map[string]interface{}, message string) {
	res := SendMessageToApi(configData, message)
	if !res {
		//进行重试
		for i := 0; i < config.RabbitmqSettingConfig["tryReConsumeNum"].(int); i++ {
			res = SendMessageToApi(configData, message)
			if res {
				break
			}
		}
	}
}