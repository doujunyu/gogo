package ff

import (
	_ "github.com/joho/godotenv/autoload"
	"net/http"
)

// HandlerFunc 接口需要执行的程序方法
type HandlerFunc func(util *Job)
type GroupFunc func()
type Centre struct {
	Middleware []HandlerFunc
}


func ReadyGo() *Centre {
	return &Centre{
		Middleware: []HandlerFunc{},
	}
}

func (c *Centre) GET(relativePath string , HandlerFunc ...HandlerFunc) {
	http.HandleFunc(relativePath, func(w http.ResponseWriter, r *http.Request) {
		job := &Job{
			Log:  JobNewLog(),  //初始化日志
			File: JobNewFile(), //初始化文件
		}
		defer func() {
			if err := recover(); err != nil {
				job.JsonError(nil, "执行错误", 500)
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




