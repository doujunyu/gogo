package ff

import (
	"gogo/utility"
)

// JsonSuccess 返回成功的json格式(data,msg,code)
func (j *Job) JsonSuccess(all ...interface{}) {
	mes := utility.Message{
		Data: make([]int, 0),
		Msg:  "操作成功",
		Code: 0,
	}
	_, _ = j.W.Write(mes.Json(all))
}

// JsonError 返回成功的json格式
func (j *Job) JsonError(all ...interface{}) {
	mes := utility.Message{
		Data: make([]int, 0),
		Msg:  "操作失败",
		Code: 1,
	}
	_, _ = j.W.Write(mes.Json(all))
}

