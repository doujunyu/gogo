package gogo

import (
	"database/sql"
	"fmt"
	"os"
	"sync"
)

var sendSqlLine *sql.DB
var sendSqlLineOnce sync.Once

// +----------------------------------------------------------------------
// | 包函数
// +----------------------------------------------------------------------

// Open 建立链接,,实现单例模式
func Open() *sql.DB {
	sendSqlLineOnce.Do(func() {
		DbConnection := os.Getenv("DB_CONNECTION")
		DbHost := os.Getenv("DB_HOST")
		DbPort := os.Getenv("DB_PORT")
		DbUserName := os.Getenv("DB_USERNAME")
		DbPassWord := os.Getenv("DB_PASSWORD")
		DbDataBase := os.Getenv("DB_DATABASE")
		DbCharset := os.Getenv("DB_CHARSET")
		sqlLine, err := sql.Open(DbConnection, DbUserName+":"+DbPassWord+"@tcp("+DbHost+":"+DbPort+")/"+DbDataBase+"?charset="+DbCharset)
		if err == nil {
			sendSqlLine = sqlLine
		}
		fmt.Println("数据库链接执行")
	})
	return sendSqlLine
}

// QueryFind 原生查询(sql语句,参数)
func QueryFind(rows *sql.Rows) ([]map[string]interface{}, error) {
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
	return list, nil
}

// Db DB方法
func Db(Table string) *Query {
	return &Query{
		RecordTable: Table,
		Tx:          nil,
	}
}

func Try() *sql.Tx {
	tx, err := Open().Begin()
	if err != nil {
		return nil
	}
	return tx
}

// +----------------------------------------------------------------------
// | 结构体接口,进行拼接数据
// +----------------------------------------------------------------------
