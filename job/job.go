package job

import (
	"net/http"
)

// Job 总工作台
type Job struct {
	W         http.ResponseWriter
	R         *http.Request
	Log       *Log                   `Testing:"日志"`
	Cache     *Cache                 `Testing:"缓存"`
	File      *Files                 `Testing:"文件"`
	Input     map[string]string      `Testing:"接收的参数"`
	GroupData map[string]interface{} `Testing:"跨方法的数据"`
	Rests     *map[string]interface{}           `Testing:"其他"`
}

// +----------------------------------------------------------------------
// | 程序执行结束，接口返回操作
// +----------------------------------------------------------------------
