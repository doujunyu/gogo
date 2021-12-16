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
	j.CheckByDataMesCode(&mes, all)
	j.W.Write(mes.Json())
}

// JsonError 返回成功的json格式
func (j *Job) JsonError(all ...interface{}) {
	mes := utility.Message{
		Data: make([]int, 0),
		Msg:  "操作失败",
		Code: 1,
	}
	j.CheckByDataMesCode(&mes, all)
	j.W.Write(mes.Json())
}

// CheckByDataMesCode 检查是否有空数据
func (j *Job) CheckByDataMesCode(mess *utility.Message, all []interface{}) {
	if len(all) >= 1 && all[0] != nil {
		mess.Data = all[0]
	}
	if len(all) >= 2 && all[1] != nil {
		mess.Msg = all[1].(string)
	}
	if len(all) >= 3 && all[2] != nil {
		mess.Code = all[2].(int)
	}
}
