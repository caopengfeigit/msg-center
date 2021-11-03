package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"runtime/debug"
)

// 错误处理的结构体
type Error struct {
	StatusCode int    `json:"-"`
	Code       int    `json:"code"`
	Msg        string `json:"message"`
	Data		interface{} `json:"data"`

}

func (e *Error) Error() string {
	return e.Msg
}

// recover错误，转string
func errorToString(r interface{}) string {
	switch v := r.(type) {
	case error:
		return v.Error()
	case string:
		return v
	default:
		return "unknown error"
	}
}

func NewError(statusCode, Code int, msg string) *Error {
	return &Error{
		StatusCode: statusCode,
		Code:       Code,
		Msg:        msg,
	}
}

// 404处理
func HandleNotFound(c *gin.Context) {
	err :=  NewError(http.StatusNotFound, 10001, http.StatusText(http.StatusNotFound))
	c.JSON(err.StatusCode,err)
	return
}
// 调用方式异常
func HandleNotAllowMethod(c *gin.Context) {
	err :=  NewError(http.StatusMethodNotAllowed, 10001, http.StatusText(http.StatusMethodNotAllowed))
	c.JSON(err.StatusCode,err)
	return
}

func ErrHandler(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			//封装通用json返回
			Err := new(Error)
			if e,ok := err.(*Error); ok {
				Err = e
			}else{
				//打印错误堆栈信息
				debug.PrintStack()
				Err.StatusCode = 500
				Err.Msg = errorToString(err)
				Err.Code = 503
			}
			c.JSON(Err.StatusCode,map[string]interface{}{"error":Err})
			//终止后续接口调用，不加的话recover到异常后，还会继续执行接口里后续代码
			c.Abort()
		}
	}()
	//加载完 defer recover，继续后续接口调用
	c.Next()
}

//抛异常方法
func AutoError(StatusCode,code int,msg string)  {
	err := new(Error)
	err.Msg = msg
	err.StatusCode = StatusCode
	err.Code = code
	panic(err)
}
//抛异常方法
func AutoErrorWithData(StatusCode,code int,msg string, data interface{})  {
	err := new(Error)
	err.Msg = msg
	err.StatusCode = StatusCode
	err.Code = code
	err.Data = data
	panic(err)
}