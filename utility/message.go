package utility

import "encoding/json"

// Message 返回数据格式
type Message struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// Json 返回json格式
func (message *Message) Json(all []interface{}) []byte {
	if len(all) >= 1 && all[0] != nil {
		message.Data = all[0]
	}
	if len(all) >= 2 && all[1] != nil {
		message.Msg = all[1].(string)
	}
	if len(all) >= 3 && all[2] != nil {
		message.Code = all[2].(int)
	}

	dataJson, err := json.Marshal(message)
	if err != nil {
		//格式不合法的json数据
		return []uint8{123, 34, 99, 111, 100, 101, 34, 58, 53, 48, 48, 44, 34, 109, 115, 103, 34, 58, 34, 230, 160, 188, 229, 188, 143, 228, 184, 141, 229, 144, 136, 230, 179, 149, 34, 44, 34, 100, 97, 116, 97, 34, 58, 91, 93, 125}
	}
	return dataJson
}
