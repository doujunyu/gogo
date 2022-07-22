package gogo

import (
	"fmt"
	"github.com/doujunyu/gogo/gogo_log"
	"github.com/doujunyu/gogo/job"
	"io/ioutil"
	"net/http"
	"runtime/debug"
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
					GroupData: make(map[string]interface{}),//初始化跨方法的数据
				}
				//接参数
				r.FormValue("")
				jobs.W, jobs.R = w, r
				if ServerStatus != ServerStatusAllow { //判断服务器是否允许访问
					jobs.JsonError(nil, "服务器停止访问...")
					return
				}
				if handlerFuncMapSlice[r.Method] == nil {
					jobs.JsonError(nil, r.Method +"请求方式不存在", 1)
					return
				}
				//参数赋值
				jobs.Input = make(map[string]string)
				for key, valuse := range r.Form {
					jobs.Input[key] = valuse[0]
				}
				//json参数赋值
				con, _ := ioutil.ReadAll(jobs.R.Body)
				defer jobs.R.Body.Close()
				jobs.InputJson = string(con)


				defer func() {
					if err := recover(); err != nil {
						logMessage := fmt.Sprintf("%v请求 路由:%v 参数%v \n 报错详情:%v \n %v",r.Method,relativePath,jobs.Input, err,string(debug.Stack()))
						gogo_log.Write("error","内部错误",  logMessage)
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












