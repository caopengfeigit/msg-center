package router

import (
	"github.com/gin-gonic/gin"
	"msgCenter/api"
	"msgCenter/controller"
	"msgCenter/handler"
)

func InitRouter(r *gin.Engine) {
	r.NoMethod(handler.HandleNotAllowMethod) //公共异常
	r.NoRoute(handler.HandleNotFound) //公共异常
	r.Use(handler.ErrHandler)
	r.StaticFile("/favicon.ico", "./public/img/favicon.ico")
	
	r.Any("/publish-message", api.RecieveMessage)

	// ========== 页面模板 start ==========
	q := r.Group("/config")
	{
		q.GET("/index", controller.Index)
		q.GET("/project", controller.Project)
		q.GET("/callback-log", controller.CallbackLog)
		q.GET("/add-config", controller.AddConfig)
	}
	// ========== 页面模板 end ==========
	
	// ========== api接口 start ==========
	apis := r.Group("/api")
	{
		apis.GET("/get-list", api.GetList)
		apis.GET("/get-project-list", api.GetProjectList)
		apis.GET("/get-config-detail", api.GetConfigDetail)
		apis.GET("/get-callback-logs", api.GetCallbackLogs)
		apis.GET("/search-project", api.SearchProject)
		apis.GET("/search-event", api.SearchEvent)
		apis.POST("/add-project", api.AddProject)
		apis.POST("/edit-project", api.EditProject)
		apis.POST("/add-config", api.AddConfig)
		apis.DELETE("/del-event-config", api.DelEventConfig)
		apis.DELETE("/del-project", api.DelProject)
	}
	// ========== api接口 end ==========
}
