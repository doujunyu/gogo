package gogo

import (
	"context"
	"fmt"
	"github.com/doujunyu/gogo/job"
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
type HandlerFunc func(util *job.Job)
type GroupFunc func()
type Centre struct {
	Middleware   []HandlerFunc `Testing:"中间件"`
	Log          *job.Log      `Testing:"日志写入管道"`
	Cache        *job.Cache    `Testing:"缓存写入管道"`
	ServerClose  chan int      `Testing:"关闭服务(传入数据执行关闭操作)"`
	ServerStatus int           `Testing:"阻止外网访问:0=正常,1=禁止,2=系统禁止(在执行关闭服务用到)"`
	Server       *http.Server  `Testing:"http服务"`
}

func ReadyGo() *Centre {
	return &Centre{
		Middleware:   []HandlerFunc{},
		Log:          job.NewLog(),
		Cache:        job.NewCache(),
		ServerClose:  make(chan int, 1),
		ServerStatus: ServerStatusAllow,
		Server: &http.Server{
			Addr: ":7070",
			Handler: http.TimeoutHandler(http.DefaultServeMux, time.Second*(60*5), func() string {
				msg := job.Message{
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

// Run 启动
func (c *Centre) Run(addr ...string) {
	go c.Cache.ChanLongTime()            //缓存
	go c.LogChanOut()                    //日志管道处理
	go c.SetClose()                      //软关闭服务
	c.Server.Addr = resolveAddress(addr) //确认端口
	_ = c.Server.ListenAndServe()        //启动
}

//缓存执行

// LogChanOut 将管道中的记录信息写入日志
func (c *Centre) LogChanOut() {
	for {
		data := <-c.Log.LogChan
		job.LogWrite(data.Url, data.FileName, data.Prefix, data.Content) //日志内存
	}
}

// SetClose 设置执行关闭服务
func (c *Centre) SetClose() {
	<-c.ServerClose
	c.ServerStatus = ServerStatusSystemForbid
	fmt.Println("http服务器已经停止外网访问!")
	close(c.Log.LogChan)
	fmt.Println("日志已停止写入!")
	fmt.Println("正在清理管道中日志信息...")
	logChanLenI := 0
	for {
		logChanLenI++
		logChanLen := len(c.Log.LogChan)
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

//根据不同的请求做出判断,私用方法
func (c *Centre) httpRequest(relativePath string, route string, HandlerFunc ...HandlerFunc) {

	http.HandleFunc(relativePath, func(w http.ResponseWriter, r *http.Request) {
		//logChan := c.LogChan
		jobs := &job.Job{
			Log:   c.Log,            //初始化日志
			Cache: c.Cache,          //换缓
			File:  job.JobNewFile(), //初始化文件
		}
		defer func() {
			if err := recover(); err != nil {
				jobs.JsonError(nil, "执行错误", 500)
				jobs.Log.Write("请求错误","error",fmt.Sprintf("%v",err))
				return
			}
		}()
		//接参数
		r.FormValue("")
		jobs.W, jobs.R = w, r
		if c.ServerStatus != ServerStatusAllow { //判断服务器是否允许访问
			jobs.JsonError(nil, "服务器停止访问...")
			return
		}
		//参数赋值
		jobs.Input = make(map[string]string)
		for key, valuse := range r.Form {
			jobs.Input[key] = valuse[0]
		}
		//判断请求方式
		if r.Method != route {
			jobs.JsonError(nil, "请求不存在", 404)
			return
		}

		//全局中间件
		for _, MiddlewareHandlerFunc := range c.Middleware {
			MiddlewareHandlerFunc(jobs)
		}
		//局部中间件
		for _, HandlerFuncVal := range HandlerFunc {
			HandlerFuncVal(jobs)
		}
	})
}

func resolveAddress(addr []string) string {
	switch len(addr) {
	case 0:
		return ":8080"
	case 1:
		return addr[0]
	default:
		panic("too many parameters")
	}
}
