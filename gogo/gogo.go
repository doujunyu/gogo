package gogo

import (
	"github.com/doujunyu/gogo/cache"
	"github.com/doujunyu/gogo/gogo_log"
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

var ServerStatus = ServerStatusAllow //阻止外网访问:0=正常,1=禁止,2=系统禁止(在执行关闭服务用到)

// HandlerFunc 接口需要执行的程序方法
type HandlerFunc func(util *job.Job)
type GroupFunc func()
type Centre struct {
	Middleware []HandlerFunc `Testing:"中间件"`
	ServerClose   chan int     `Testing:"关闭服务(传入数据执行关闭操作)"`
	Server        *http.Server `Testing:"http服务"`
	gatherRequest map[string]map[string]*[]HandlerFunc `Testing:"路由接口集合"`
}

func ReadyGo() *Centre {
	return &Centre{
		Middleware:   []HandlerFunc{},
		ServerClose:  make(chan int, 1),
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
		gatherRequest: make(map[string]map[string]*[]HandlerFunc),
	}
}

// Run 启动
func (c *Centre) Run(addr ...interface{}) {
	c.createRequestMapDataRun()          //生成路由接口
	go c.setClose()                      //软关闭服务
	go gogo_log.LogChanOut()             //日志管道处理
	go cache.ChanLongTime()              //缓存清除过期数据
	c.Server.Addr = resolveAddress(addr) //确认端口
	_ = c.Server.ListenAndServe()        //启动
}



