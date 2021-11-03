package base

import (
	"github.com/gin-gonic/gin"
	"msgCenter/config"
)

//获取参数
func GetParams(c *gin.Context) (requestParams config.RequestParams, valid bool) {
	valid = true
	if err := c.BindQuery(&requestParams); err != nil {
		valid = false
	}
	if c.Request.Method == "POST" {
		if err := c.ShouldBind(&requestParams); err != nil {
			valid = false
		}
	}
	//兼容老系统，在这进行数据校验
	if requestParams.Message == "" || requestParams.ConfigKey == "" {
		valid = false
	}
	return requestParams, valid
}