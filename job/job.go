package job

import (
	"net/http"
)

// Job 总工作台
type Job struct {
	W     http.ResponseWriter
	R     *http.Request
	Log   *Log                      //日志
	Cache *Cache
	File  *Files                    //文件
	Input map[string]string         //接收的参数
	GroupData map[string]interface{} //跨方法的数据
}

// +----------------------------------------------------------------------
// | 程序执行结束，接口返回操作
// +----------------------------------------------------------------------
