package ff

import (
	"fmt"
	_ "github.com/joho/godotenv/autoload"
	"gogo/utility"
	"net/http"
	"time"
)

// HandlerFunc 接口需要执行的程序方法
type HandlerFunc func(util *Job)
type GroupFunc func()
type Centre struct {
	Middleware []HandlerFunc
	LogChan *chan utility.LogWriteStrings
	Server *http.Server
}


func ReadyGo() *Centre {
	logChan := make(chan utility.LogWriteStrings)

	return &Centre{
		Middleware: []HandlerFunc{},
		LogChan: &logChan,
		Server: &http.Server{
			Addr:           ":7070",
			Handler:        http.TimeoutHandler(http.DefaultServeMux, time.Second * (60 * 5), func()string{
				msg := utility.Message{
					Data: make([]int, 0),
					Msg:  "操作失败",
					Code: 1,
				}
				return string(msg.Json(nil))
			}()),
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
			MaxHeaderBytes: 1 << 20,
		},
	}
}
func (c *Centre) GET(relativePath string , HandlerFunc ...HandlerFunc) {

	http.HandleFunc(relativePath, func(w http.ResponseWriter, r *http.Request) {
		job := &Job{
			Log:  JobNewLog(c.LogChan),  //初始化日志
			File: JobNewFile(), //初始化文件
		}
		defer func() {
			if err := recover(); err != nil {
				job.JsonError(nil, "执行错误", 500)
				fmt.Println(err)
				return
			}
		}()
		//接参数
		r.FormValue("")
		job.W, job.R = w, r
		//参数赋值
		job.Input = make(map[string]string)
		for key, valuse := range r.Form {
			job.Input[key] = valuse[0]
		}
		//判断请求方式
		if r.Method != "GET" {
			job.JsonError(nil, "请求不存在", 404)
			return
		}
		//全局中间件
		for _, MiddlewareHandlerFunc := range c.Middleware {
			MiddlewareHandlerFunc(job)
		}
		//局部中间件
		for _, HandlerFuncVal := range HandlerFunc {
			HandlerFuncVal(job)
		}
		//sql.Exit()
	})
}

func (c *Centre) LogChanOut(){
	for{
		data := <- *c.LogChan
		utility.LogWrite(data.Url, data.FileName, data.Prefix, data.Content) //日志内存
	}
}




