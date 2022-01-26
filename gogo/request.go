package gogo

import (
	"context"
	"fmt"
	"github.com/doujunyu/gogo/gogo_log"
	"github.com/doujunyu/gogo/job"
	"net/http"
	"runtime/debug"
	"time"
)

// GET 请求
func (c *Centre) GET(relativePath string, handlerFunc ...HandlerFunc) {
	c.requestMapData(relativePath, "GET", handlerFunc...)
}
// POST 请求
func (c *Centre) POST(relativePath string, HandlerFunc ...HandlerFunc) {
	c.requestMapData(relativePath, "POST", HandlerFunc...)
}
// PUT 请求
func (c *Centre) PUT(relativePath string, HandlerFunc ...HandlerFunc) {
	c.requestMapData(relativePath, "POST", HandlerFunc...)
}
// PATCH 请求
func (c *Centre) PATCH(relativePath string, HandlerFunc ...HandlerFunc) {
	c.requestMapData(relativePath, "POST", HandlerFunc...)
}
// DELETE 请求
func (c *Centre) DELETE(relativePath string, HandlerFunc ...HandlerFunc) {
	c.requestMapData(relativePath, "POST", HandlerFunc...)
}
// COPY 请求
func (c *Centre) COPY(relativePath string, HandlerFunc ...HandlerFunc) {
	c.requestMapData(relativePath, "POST", HandlerFunc...)
}
// OPTIONS 请求
func (c *Centre) OPTIONS(relativePath string, HandlerFunc ...HandlerFunc) {
	c.requestMapData(relativePath, "POST", HandlerFunc...)
}
// LINK 请求
func (c *Centre) LINK(relativePath string, HandlerFunc ...HandlerFunc) {
	c.requestMapData(relativePath, "POST", HandlerFunc...)
}
// UNLINK 请求
func (c *Centre) UNLINK(relativePath string, HandlerFunc ...HandlerFunc) {
	c.requestMapData(relativePath, "POST", HandlerFunc...)
}
// PURGE 请求
func (c *Centre) PURGE(relativePath string, HandlerFunc ...HandlerFunc) {
	c.requestMapData(relativePath, "POST", HandlerFunc...)
}
// LOCK 请求
func (c *Centre) LOCK(relativePath string, HandlerFunc ...HandlerFunc) {
	c.requestMapData(relativePath, "POST", HandlerFunc...)
}
// UNLOCK 请求
func (c *Centre) UNLOCK(relativePath string, HandlerFunc ...HandlerFunc) {
	c.requestMapData(relativePath, "POST", HandlerFunc...)
}
// PROPFIND 请求
func (c *Centre) PROPFIND(relativePath string, HandlerFunc ...HandlerFunc) {
	c.requestMapData(relativePath, "POST", HandlerFunc...)
}
// VIEW 请求
func (c *Centre) VIEW(relativePath string, HandlerFunc ...HandlerFunc) {
	c.requestMapData(relativePath, "POST", HandlerFunc...)
}



// SetClose 设置执行关闭服务
func (c *Centre) setClose() {
	<-c.ServerClose
	ServerStatus = ServerStatusSystemForbid
	fmt.Println("http服务器已经停止外网访问!")
	close(gogo_log.GlobalLogData.LogChan)
	fmt.Println("日志已停止写入!")
	fmt.Println("正在清理管道中日志信息...")
	logChanLenI := 0
	for {
		logChanLenI++
		logChanLen := len(gogo_log.GlobalLogData.LogChan)
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

// 组装所有接口路径
func (c *Centre) requestMapData(relativePath string, route string, handlerFunc ...HandlerFunc) {
	if c.gatherRequest[relativePath] == nil {
		request := make(map[string]*[]HandlerFunc)
		request[route] = &handlerFunc
		c.gatherRequest[relativePath] = request
	} else {
		c.gatherRequest[relativePath][route] = &handlerFunc
	}

}

//根据不同的请求做出判断,私用方法
func (c *Centre) createRequestMapDataRun() {
	for key, val := range c.gatherRequest {
		func(relativePath string, handlerFuncMapSlice map[string]*[]HandlerFunc) {
			http.HandleFunc(relativePath, func(w http.ResponseWriter, r *http.Request) {

				//logChan := c.LogChan
				jobs := &job.Job{
					File:  job.JobNewFile(), //初始化文件
					IsFlow: true,
				}
				//接参数
				r.FormValue("")
				jobs.W, jobs.R = w, r
				if ServerStatus != ServerStatusAllow { //判断服务器是否允许访问
					jobs.JsonError(nil, "服务器停止访问...")
					return
				}
				//参数赋值
				jobs.Input = make(map[string]string)
				for key, valuse := range r.Form {
					jobs.Input[key] = valuse[0]
				}

				if handlerFuncMapSlice[r.Method] == nil {
					jobs.JsonError(nil, "请求方式不存在", 1)
					return
				}
				defer func() {
					if err := recover(); err != nil {
						logMessage := fmt.Sprintf("%v请求 路由:%v 参数%v \n 报错详情:%v \n %v",r.Method,relativePath,jobs.Input, err,string(debug.Stack()))
						gogo_log.Write("内部错误", "error", logMessage)
						jobs.JsonError(nil, "执行错误", 500)
						return
					}
				}()
				//全局中间件
				for _, MiddlewareHandlerFunc := range c.Middleware {
					if jobs.IsFlow == true {
						MiddlewareHandlerFunc(jobs)
					}
				}
				//局部中间件
				for _, HandlerFuncVal := range *handlerFuncMapSlice[r.Method] {
					if jobs.IsFlow == true {
						HandlerFuncVal(jobs)
					}
				}
			})
		}(key, val)
	}

}

func resolveAddress(addr []interface{}) string {
	switch len(addr) {
	case 0:
		return ":8080"
	case 1:
		return addr[0].(string)
	default:
		panic("too many parameters")
	}
}











