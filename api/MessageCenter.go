package api

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"msgCenter/base"
	"msgCenter/config"
	"msgCenter/db"
	"msgCenter/logic"
	"msgCenter/model"
	"net/http"
	"strconv"
	"strings"
	"time"
)

//接收异步消息（消息的发布者）
func RecieveMessage(c *gin.Context) {
	requestParams, valid := base.GetParams(c)
	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{"status" : 10001, "error": config.ErrorMessage[10001]})
		return
	}
	//移交给logic进行处理
	code := logic.HandleMessage(requestParams)
	//释放rabbitmq连接资源
	db.CloseRabbitmqResource()
	c.Header("Content-type", "application/json")
	if code > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"status" : code, "error": config.ErrorMessage[code]})
	} else if code == 0 {
		c.JSON(http.StatusOK, gin.H{"status" : code, "error" : ""})
	}
}

/**
 * 获取配置列表
 */
func GetList(c *gin.Context) {
	projectName := c.DefaultQuery("ProjectName", "")
	eventName := c.DefaultQuery("EventName", "")
	page := c.DefaultQuery("page", "1")
	//页码转换
	convNum, err := strconv.Atoi(page)
	num := 1
	if err == nil {
		num = convNum
	}
	list, pagination := logic.GetList(projectName, eventName, num)
	c.JSON(http.StatusOK, gin.H{"list" : list, "pagination" : pagination})
}

//获取业务线列表
func GetProjectList(c *gin.Context) {
	name := c.DefaultQuery("Name", "")
	businessLine := c.DefaultQuery("BusinessLine", "")
	page := c.DefaultQuery("page", "1")
	//页码转换
	convNum, err := strconv.Atoi(page)
	num := 1
	if err == nil {
		num = convNum
	}
	list, pagination := logic.GetProjectList(name, businessLine, num)
	c.JSON(http.StatusOK, gin.H{"list" : list, "pagination" : pagination})
}

//删除事件及对应的数据
func DelEventConfig(c *gin.Context) {
	//要删除的事件配置id
	id := c.DefaultQuery("id", "")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H {
			"message" : "请选择删除的配置",
		})
		return
	}
	id = strings.Trim(id, ",")
	success, code := logic.DelEventConfig(id)
	if !success {
		c.JSON(http.StatusBadRequest, gin.H{
			"message" : config.ErrorMessage[code],
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"message" : "删除成功",
		})
	}
}

//删除业务线及对应的数据
func DelProject(c *gin.Context) {
	//要删除的业务线id
	id := c.DefaultQuery("id", "")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H {
			"message" : "请选择删除的业务线",
		})
		return
	}
	id = strings.Trim(id, ",")
	success, code := logic.DelProject(id)
	if !success {
		c.JSON(http.StatusBadRequest, gin.H{
			"message" : config.ErrorMessage[code],
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"message" : "删除成功",
		})
	}
}

//添加业务线
func AddProject(c *gin.Context) {
	//标识和业务线
	name := c.DefaultPostForm("name", "")
	businessLine := c.DefaultPostForm("business_line", "")
	if name == "" || businessLine == "" {
		c.JSON(http.StatusBadRequest, gin.H {
			"message" : "请填写必填项",
		})
		return
	}
	success, code := logic.AddProject(name, businessLine)
	if !success {
		c.JSON(http.StatusBadRequest, gin.H{
			"message" : config.ErrorMessage[code],
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"message" : "添加成功",
		})
	}
}

//修改业务线
func EditProject(c *gin.Context) {
	//id 标识 业务线
	id := c.DefaultPostForm("id", "")
	name := c.DefaultPostForm("name", "")
	businessLine := c.DefaultPostForm("business_line", "")
	if id == "" || name == "" || businessLine == "" {
		c.JSON(http.StatusBadRequest, gin.H {
			"message" : "请填写必填项",
		})
		return
	}
	//转换id
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H {
			"message" : "参数错误",
		})
		return
	}
	
	success, code := logic.EditProject(objectId, name, businessLine)
	if !success {
		c.JSON(http.StatusBadRequest, gin.H{
			"message" : config.ErrorMessage[code],
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"message" : "修改成功",
		})
	}
}

//获取事件配置详情
func GetConfigDetail(c *gin.Context) {
	//eventid
	id := c.DefaultQuery("id", "")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H {
			"message" : "参数错误",
		})
		return
	}
	//转换id
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H {
			"message" : "参数错误",
		})
		return
	}
	
	configDetail, code := logic.GetConfigDetail(objectId)
	if code > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message" : config.ErrorMessage[code],
		})
	} else {
		c.JSON(http.StatusOK, configDetail)
	}
}

//获取回调记录列表
func GetCallbackLogs(c *gin.Context) {
	projectName := c.DefaultQuery("ProjectName", "")
	eventName := c.DefaultQuery("EventKey", "")
	startDate := c.DefaultQuery("StartDate", "")
	endDate := c.DefaultQuery("EndDate", "")
	page := c.DefaultQuery("page", "1")
	var startDateConv, endDateConv time.Time
	//页码转换
	convNum, err := strconv.Atoi(page)
	num := 1
	if err == nil {
		num = convNum
	}
	//时间转换时区
	if len(startDate) > 0 {
		if startDate > endDate {
			c.JSON(http.StatusBadRequest, gin.H {
				"message" : "时间参数错误",
			})
			return
		}
		var cstZone = time.FixedZone("CST", 8*3600)
		startDateConv, _ = time.ParseInLocation("2006-01-02T15:04:05Z", startDate, cstZone)
		endDateConv, _ = time.ParseInLocation("2006-01-02T15:04:05Z", endDate, cstZone)
	}
	list, pagination := logic.GetCallbackLogs(projectName, eventName, startDateConv, endDateConv, num)
	c.JSON(http.StatusOK, gin.H{"list" : list, "pagination" : pagination})
}

//搜索业务线
func SearchProject(c *gin.Context) {
	projectName := c.DefaultQuery("projectName", "")
	var project model.Project
	var hasProject bool
	if len(projectName) > 0 {
		project = logic.SearchProject(projectName)
		if !project.Id.IsZero() {
			hasProject = true
		}
	}
	c.JSON(http.StatusOK, gin.H{"hasProject" : hasProject, "project" : project})
}

//搜索事件是否存在
func SearchEvent(c *gin.Context) {
	projectId := c.DefaultQuery("projectId", "")
	eventName := c.DefaultQuery("eventName", "")
	if len(projectId) == 0 || len(eventName) == 0 {
		c.JSON(http.StatusBadRequest, gin.H {
			"message" : "参数错误",
		})
		return
	}
	//转换id
	objectId, err := primitive.ObjectIDFromHex(projectId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H {
			"message" : "参数错误",
		})
		return
	}
	var hasEvent bool
	event := logic.SearchEvent(objectId, eventName)
	if !event.Id.IsZero() {
		hasEvent = true
	}
	c.JSON(http.StatusOK, gin.H{"hasEvent" : hasEvent})
}

//添加配置
func AddConfig(c *gin.Context) {
	var params config.AddConfigRequestParams
	if err := c.ShouldBind(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H {
			"message" : "参数错误",
		})
		return
	}
	success, code := logic.AddConfig(params)
	if !success {
		c.JSON(http.StatusBadRequest, gin.H{
			"message" : config.ErrorMessage[code],
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"message" : "添加成功",
		})
	}
}