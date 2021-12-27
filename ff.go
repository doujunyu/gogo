package gogo

import (
	"context"
	"fmt"
	_ "github.com/joho/godotenv/autoload"
	"net/http"
	"time"
)

const (
	ServerStatusAllow        = 0 //正常
	ServerStatusForbid       = 1 //禁止
	ServerStatusSystemForbid = 2 //系统禁止
)

// HandlerFunc 接口需要执行的程序方法
type HandlerFunc func(util *Job)
type GroupFunc func()
type Centre struct {
	Middleware   []HandlerFunc         `Testing:"中间件"`
	LogChan      *chan LogWriteStrings `Testing:"日志写入管道"`
	Server       *http.Server          `Testing:"http服务"`
	ServerClose  *chan int             `Testing:"关闭服务(传入数据执行关闭操作)"`
	ServerStatus *int                  `Testing:"阻止外网访问:0=正常,1=禁止,2=系统禁止(在执行关闭服务用到)"`
}

func ReadyGo() *Centre {
	logChan := make(chan LogWriteStrings, 1000)
	serverClose := make(chan int, 1)
	serverStatus := ServerStatusAllow
	return &Centre{
		Middleware:   []HandlerFunc{},
		LogChan:      &logChan,
		ServerClose:  &serverClose,
		ServerStatus: &serverStatus,
		Server: &http.Server{
			Addr: ":7070",
			Handler: http.TimeoutHandler(http.DefaultServeMux, time.Second*(60*5), func() string {
				msg := Message{
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

// GET 请求
func (c *Centre) GET(relativePath string, HandlerFunc ...HandlerFunc) {
	c.httpRequest(relativePath, "GET", HandlerFunc...)
}

// POST 请求
func (c *Centre) POST(relativePath string, HandlerFunc ...HandlerFunc) {
	c.httpRequest(relativePath, "POST", HandlerFunc...)
}

// PUT 请求
func (c *Centre) PUT(relativePath string, HandlerFunc ...HandlerFunc) {
	c.httpRequest(relativePath, "POST", HandlerFunc...)
}

// PATCH 请求
func (c *Centre) PATCH(relativePath string, HandlerFunc ...HandlerFunc) {
	c.httpRequest(relativePath, "POST", HandlerFunc...)
}

// DELETE 请求
func (c *Centre) DELETE(relativePath string, HandlerFunc ...HandlerFunc) {
	c.httpRequest(relativePath, "POST", HandlerFunc...)
}

// COPY 请求
func (c *Centre) COPY(relativePath string, HandlerFunc ...HandlerFunc) {
	c.httpRequest(relativePath, "POST", HandlerFunc...)
}

// OPTIONS 请求
func (c *Centre) OPTIONS(relativePath string, HandlerFunc ...HandlerFunc) {
	c.httpRequest(relativePath, "POST", HandlerFunc...)
}

// LINK 请求
func (c *Centre) LINK(relativePath string, HandlerFunc ...HandlerFunc) {
	c.httpRequest(relativePath, "POST", HandlerFunc...)
}

// UNLINK 请求
func (c *Centre) UNLINK(relativePath string, HandlerFunc ...HandlerFunc) {
	c.httpRequest(relativePath, "POST", HandlerFunc...)
}

// PURGE 请求
func (c *Centre) PURGE(relativePath string, HandlerFunc ...HandlerFunc) {
	c.httpRequest(relativePath, "POST", HandlerFunc...)
}

// LOCK 请求
func (c *Centre) LOCK(relativePath string, HandlerFunc ...HandlerFunc) {
	c.httpRequest(relativePath, "POST", HandlerFunc...)
}

// UNLOCK 请求
func (c *Centre) UNLOCK(relativePath string, HandlerFunc ...HandlerFunc) {
	c.httpRequest(relativePath, "POST", HandlerFunc...)
}

// PROPFIND 请求
func (c *Centre) PROPFIND(relativePath string, HandlerFunc ...HandlerFunc) {
	c.httpRequest(relativePath, "POST", HandlerFunc...)
}

// VIEW 请求
func (c *Centre) VIEW(relativePath string, HandlerFunc ...HandlerFunc) {
	c.httpRequest(relativePath, "POST", HandlerFunc...)
}

// SetClose 设置执行关闭服务
func (c *Centre) SetClose() {
	<-*c.ServerClose
	*c.ServerStatus = ServerStatusSystemForbid
	fmt.Println("http服务器已经停止外网访问!")
	close(*c.LogChan)
	fmt.Println("日志已停止写入!")
	fmt.Println("正在清理管道中日志信息...")
	logChanLenI := 0
	for {
		logChanLenI++
		logChanLen := len(*c.LogChan)
		if logChanLen == 0 {
			break
		}
		fmt.Print(logChanLen, "->")
		time.Sleep(time.Second)
		if logChanLenI == 10 {
			logChanLenI = 0
			fmt.Println()
		}
	}
	fmt.Println("日志清理完毕")
	fmt.Println("15秒后关闭计算机...")
	tx, cancel := context.WithTimeout(context.TODO(), 20*time.Second)
	defer cancel()
	for i := 15; i > 0; i-- {
		if i == 0 || i == 5 || i == 10 {
			fmt.Println()
		}
		time.Sleep(time.Second)
		fmt.Print(i, "->")
	}
	fmt.Println("正在关闭...")
	_ = c.Server.Shutdown(tx)
}

// LogChanOut 将管道中的记录信息写入日志
func (c *Centre) LogChanOut() {
	for {
		data := <-*c.LogChan
		LogWrite(data.Url, data.FileName, data.Prefix, data.Content) //日志内存
	}
}

//根据不同的请求做出判断,私用方法
func (c *Centre) httpRequest(relativePath string, route string, HandlerFunc ...HandlerFunc) {

	http.HandleFunc(relativePath, func(w http.ResponseWriter, r *http.Request) {
		job := &Job{
			Log:  JobNewLog(c.LogChan), //初始化日志
			File: JobNewFile(),         //初始化文件
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
		if *c.ServerStatus != ServerStatusAllow { //判断服务器是否允许访问
			job.JsonError(nil, "服务器停止访问...")
			return
		}
		//参数赋值
		job.Input = make(map[string]string)
		for key, valuse := range r.Form {
			job.Input[key] = valuse[0]
		}
		//判断请求方式
		if r.Method != route {
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
