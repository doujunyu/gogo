package job

import (
	"net/http"
)

// Job 总工作台
type Job struct {
	W         http.ResponseWriter
	R         *http.Request
	File      *Files                 `Testing:"文件"`
	Input     map[string]string      `Testing:"接收的参数"`
	InputJson string                 `Testing:"接收的json数据"`
	GroupData map[string]interface{} `Testing:"跨方法的数据"`
	IsFlow    bool                   `Testing:"判断是否继续向下执行:true=继续执行,false=停止"`
}

// +----------------------------------------------------------------------
// | 程序执行结束，接口返回操作
// +----------------------------------------------------------------------
