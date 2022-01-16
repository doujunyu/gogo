package sql_aid

import (
	"database/sql"
)


// +----------------------------------------------------------------------
// | 包函数
// +----------------------------------------------------------------------

// Open 建立链接,,实现单例模式
func Open(sqlType string,databaseLine string) (*sql.DB,error) {
	sendMySqlLine, err := sql.Open(sqlType, databaseLine)
	if err != nil {
		return nil,err
	}
	return sendMySqlLine,nil
}

// DataToMap 处理数据
func DataToMap(rows *sql.Rows,err error) ([]map[string]interface{}, error) {
	if err != nil{
		return nil,err
	}
	columns, _ := rows.Columns() //数据的字段
	columnLength := len(columns)
	cache := make([]interface{}, columnLength) //临时存储每行数据
	for index, _ := range cache {              //为每一列初始化一个指针
		var a interface{}
		cache[index] = &a
	}
	var list []map[string]interface{} //返回的切片
	for rows.Next() {                 //循环每一条数据
		_ = rows.Scan(cache...) //给每一条数据上的字段赋值
		item := make(map[string]interface{})
		for i, data := range cache { //循环每个字段
			//fmt.Printf("当前字段：%v 字段类型：%T \n",columns[i],*data.(*interface{}))
			switch (*data.(*interface{})).(type) { //判断每个字段的类型并处理
			case nil:
				item[columns[i]] = nil
			case int64:
				item[columns[i]] = *data.(*interface{}) //本身是int类型
			case []uint8:
				item[columns[i]] = string((*data.(*interface{})).([]byte)) //转成字符串
			default:
				item[columns[i]] = *data.(*interface{}) //不知道的类型不动
			}
		}
		list = append(list, item)
	}
	_ = rows.Close()
	return list, nil
}
