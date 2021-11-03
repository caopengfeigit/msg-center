package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

//配置列表
func Index(c *gin.Context) {
	c.HTML(http.StatusOK, "default/index.html", gin.H{})
}

//业务线列表
func Project(c *gin.Context) {
	c.HTML(http.StatusOK, "default/project.html", gin.H{})
}

//回调记录列表
func CallbackLog(c *gin.Context) {
	c.HTML(http.StatusOK, "default/callback.html", gin.H{})
}

//添加配置
func AddConfig(c *gin.Context) {
	c.HTML(http.StatusOK, "default/addConfig.html", gin.H{})
}