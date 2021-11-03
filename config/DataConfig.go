package config

import "go.mongodb.org/mongo-driver/bson/primitive"

var CloseChan = make(chan primitive.ObjectID)

//每页数据获取条数
var PageSize int = 20

//分页数据
type Pagination struct {
	//页码
	Page int
	//每页条数
	PageSize int
	//总页数
	TotalPage int
	//数据总条数
	Total int
}