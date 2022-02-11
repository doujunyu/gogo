package sql_aid

import (
	"fmt"
	"github.com/doujunyu/gogo/utility"
	"reflect"
	"strconv"
	"strings"
)

type PgQuery struct {
	RecordTable   string        `json:"record_table" Testing:"表明"`
	RecordField   []string      `json:"record_field" Testing:"字段"`
	RecordOrder   []string        `json:"record_order" Testing:"排序"`
	RecordGroup   []string        `json:"record_group" Testing:"分组"`
	RecordPage    string           `json:"record_page" Testing:"页数"`
	RecordSize    string           `json:"record_size" Testing:"每页数据量"`
	SqlQuery      string        `json:"sql_query,string" Testing:"sql语句"`
	WhereSqlQuery string        `json:"where_sql_query" Testing:"sql条件"`
	Args          []interface{} `json:"args" Testing:"值"`

}
type PgChildQuery func(*PgQuery, ...interface{})
// +----------------------------------------------------------------------
// | 查询
// +----------------------------------------------------------------------

//查询数据方法

func (db *PgQuery) ToSql() (string, []interface{}) {
	db.jointSql()
	db.replacePlace()
	return db.SqlQuery, db.Args
}
func (db *PgQuery) jointSql() (string, []interface{}) {
	db.OperateFindToSql()
	return db.SqlQuery, db.Args
}
//查询固定方法

func (db *PgQuery) Table(Table string) *PgQuery {
	db.RecordTable = Table
	return db
}
func (db *PgQuery) Field(field ...string) *PgQuery {
	for key, val := range field {
		field[key] = val
	}
	db.RecordField = field
	return db
}
func (db *PgQuery) OrderBy(Order string) *PgQuery {
	db.RecordOrder = append(db.RecordOrder,Order)
	return db
}
func (db *PgQuery) GroupBy(groupBy string) *PgQuery {
	db.RecordGroup = append(db.RecordGroup,groupBy)
	return db
}

//where条件

func (db *PgQuery) Where(field string, val interface{}) *PgQuery {
	if db.WhereSqlQuery != "" {
		db.WhereSqlQuery += "and "
	}
	db.WhereSqlQuery += field + " "
	if val != nil{
		db.Args = append(db.Args, val)
	}
	return db
}
func (db *PgQuery) WhereOr(field string, val interface{}) *PgQuery {
	if db.WhereSqlQuery != "" {
		db.WhereSqlQuery += "and "
	}else{
		db.WhereSqlQuery += "OR "
	}
	db.WhereSqlQuery += "(" + field + ") "
	if val != nil{
		db.Args = append(db.Args, val)
	}
	return db
}
func (db *PgQuery) WhereIn(field string, condition ...interface{}) *PgQuery {
	if db.WhereSqlQuery != "" {
		db.WhereSqlQuery += "and "
	}
	db.WhereSqlQuery += field + " in ("
	for _, _ = range condition {
		db.WhereSqlQuery += "?,"
	}
	db.WhereSqlQuery = db.WhereSqlQuery[:len(db.WhereSqlQuery)-1]
	db.WhereSqlQuery += ") "
	db.Args = append(db.Args, condition...)

	return db
}
func (db *PgQuery) WhereNotIn(field string, condition ...interface{}) *PgQuery {
	if db.WhereSqlQuery != "" {
		db.WhereSqlQuery += "and "
	}
	db.WhereSqlQuery += field + " not in ("
	for _, _ = range condition {
		db.WhereSqlQuery += "?,"
	}
	db.WhereSqlQuery = db.WhereSqlQuery[:len(db.WhereSqlQuery)-1]
	db.WhereSqlQuery += ") "
	db.Args = append(db.Args, condition...)

	return db
}
func (db *PgQuery) WhereRaw(childQuery PgChildQuery, val ...interface{}) *PgQuery {
	check := &PgQuery{}
	childQuery(check, val...)
	checkSql,args := check.jointSql()
	if checkSql == "" {
		return db
	}
	if db.WhereSqlQuery != "" {
		db.WhereSqlQuery += "and "
	}
	db.WhereSqlQuery += "("
	db.WhereSqlQuery += checkSql
	db.Args = append(db.Args, args...)
	db.WhereSqlQuery += ") "
	return db
}
func (db *PgQuery) WhereOrRaw(childQuery PgChildQuery, val ...interface{}) *PgQuery {
	check := &PgQuery{}
	childQuery(check, val...)
	checkSql,args := check.jointSql()
	if checkSql == "" {
		return db
	}
	if db.WhereSqlQuery != "" {
		db.WhereSqlQuery += "OR ("
	}else{
		db.WhereSqlQuery += "( "
	}
	db.WhereSqlQuery += checkSql
	db.Args = append(db.Args, args...)
	db.WhereSqlQuery += ") "
	return db
}
func (db *PgQuery) WhereInRaw(field string, childQuery PgChildQuery, val ...interface{}) *PgQuery {
	if db.WhereSqlQuery != "" {
		db.WhereSqlQuery += "and "
	}
	db.WhereSqlQuery += field + " in ("
	check := &PgQuery{}
	childQuery(check, val...)
	checkSql,args := check.jointSql()
	db.WhereSqlQuery += checkSql
	db.Args = append(db.Args, args...)
	db.WhereSqlQuery += ") "
	return db
}
func (db *PgQuery) WhereNotInRaw(field string, childQuery PgChildQuery, val ...interface{}) *PgQuery {
	if db.WhereSqlQuery != "" {
		db.WhereSqlQuery += "and "
	}
	db.WhereSqlQuery += field + "not in ("
	check := &PgQuery{}
	childQuery(check, val...)
	checkSql,args := check.jointSql()
	db.WhereSqlQuery += checkSql
	db.Args = append(db.Args, args...)
	db.WhereSqlQuery += ") "
	return db
}
func (db *PgQuery) WhereId(id interface{}) *PgQuery {
	if db.WhereSqlQuery != "" {
		db.WhereSqlQuery += "and "
	}
	db.WhereSqlQuery = "id = ? "
	db.Args = append(db.Args, id)
	return db
}
func (db *PgQuery) PageSize(page string, size string) *PgQuery {
	db.RecordPage = page
	db.RecordSize = size
	return db
}

// 整理查询的sql和参数

func (db *PgQuery) OperateFindToSql() {
	if db.RecordTable != "" {
		db.OperateFindField()
		db.OperateFindTable()
	}
	if db.RecordTable != "" && db.WhereSqlQuery != "" {
		db.SqlQuery += "where "
	}
	db.SqlQuery += db.WhereSqlQuery
	db.OperateFindGroupBy()
	db.OperateFindOrderBy()
	db.OperateFindPageSize()
}
func (db *PgQuery) OperateFindField() {
	db.SqlQuery = "select "
	if db.RecordField != nil {
		db.SqlQuery += utility.StringBySliceString(",", db.RecordField) + " "
	} else {
		db.SqlQuery += "* "
	}
}
func (db *PgQuery) OperateFindTable() {
	db.SqlQuery += "FROM " + db.RecordTable + " "
}
func (db *PgQuery) OperateFindGroupBy() {
	if db.RecordGroup != nil{
		db.SqlQuery += "GROUP BY " + strings.Join(db.RecordGroup,",") + " "
	}
}
func (db *PgQuery) OperateFindOrderBy() {
	if db.RecordOrder != nil{
		db.SqlQuery += "ORDER BY " + strings.Join(db.RecordOrder,",") + " "
	}
}

func (db *PgQuery) OperateFindPageSize() {
	page,_ := strconv.Atoi(db.RecordPage)
	if page != 0 {
		size,_ := strconv.Atoi(db.RecordSize)
		if size == 0 {
			size = 10
		}
		db.SqlQuery += "limit ? OFFSET ? "
		db.Args = append(db.Args, size)
		db.Args = append(db.Args, (page - 1) * size)
	}

}

// +----------------------------------------------------------------------
// | 添加方法
// +----------------------------------------------------------------------

func (db *PgQuery) InsertByMap(data *map[string]interface{}) (string,[]interface{}) {
	db.OperateInsertTable()
	db.OperateInsertDataByMap(data)
	db.replacePlace()
	return db.SqlQuery, db.Args
}
func (db *PgQuery) InsertByStruct(data interface{}) (string,[]interface{}) {
	db.OperateInsertTable()
	db.OperateInsertDataByStruct(data)
	db.replacePlace()
	return db.SqlQuery, db.Args
}
func (db *PgQuery) InsertAllByMap(datas *[]map[string]interface{}) (string,[]interface{}) {
	db.OperateInsertTable()
	for key, val := range *datas {
		if key == 0 {
			db.OperateInsertDataByMap(&val)
		} else {
			db.OperateInsertDataByMapValue(&val)
		}
	}
	db.replacePlace()
	return db.SqlQuery, db.Args

}
func (db *PgQuery) InsertAllByStruct(datas []interface{}) (string,[]interface{}) {
	db.OperateInsertTable()
	for key, val := range datas {
		if key == 0 {
			db.OperateInsertDataByStruct(val)
		} else {
			db.OperateInsertDataByStructValue(val)
		}
	}
	db.replacePlace()
	return db.SqlQuery, db.Args


}

// 整理查询的sql和参数

func (db *PgQuery) OperateInsertTable() {
	if db.RecordTable != "" {
		db.SqlQuery += "INSERT INTO " + db.RecordTable + " "
	}
}
func (db *PgQuery) OperateInsertDataByMap(data *map[string]interface{}) {
	numData := len(*data)
	if numData > 0 {
		db.SqlQuery += "("
		values := ""
		for key, val := range *data {
			db.SqlQuery +=  key + ","
			values += "?,"
			db.Args = append(db.Args, val)
		}
		db.SqlQuery = db.SqlQuery[:len(db.SqlQuery)-1]
		values = values[:len(values)-1]
		db.SqlQuery += ") VALUES (" + values + ")"
	}

}
func (db *PgQuery) OperateInsertDataByMapValue(data *map[string]interface{}) {
	db.SqlQuery += ",( "
	values := ""
	numData := len(*data)
	if numData > 0 {
		for _, val := range *data {
			values += "?,"
			db.Args = append(db.Args, val)
		}
		values = values[:len(values)-1]
		db.SqlQuery += values + ")"
	}
}
func (db *PgQuery) OperateInsertDataByStruct(data interface{}) {
	dataType := reflect.TypeOf(data).Elem()
	dataValue := reflect.ValueOf(data).Elem()
	if dataType.Kind() != reflect.Struct {
		return
	}
	numField := dataValue.NumField()
	if numField > 0 {
		values := ""
		db.SqlQuery += "("
		for i := 0; i < dataValue.NumField(); i++ {
			field := dataType.Field(i).Tag.Get("json")
			structValue := dataValue.Field(i).Interface()
			db.SqlQuery += field + ","
			values += "?,"
			db.Args = append(db.Args, structValue)
		}
		db.SqlQuery = db.SqlQuery[:len(db.SqlQuery)-1]
		values = values[:len(values)-1]
		db.SqlQuery += ") VALUES (" + values + ") "
	}
}
func (db *PgQuery) OperateInsertDataByStructValue(data interface{}) {
	dataType := reflect.TypeOf(data).Elem()
	dataValue := reflect.ValueOf(data).Elem()
	if dataType.Kind() != reflect.Struct {
		return
	}
	numField := dataValue.NumField()
	values := ",("
	if numField > 0 {
		for i := 0; i < numField; i++ {

			structValue := dataValue.Field(i).Interface()

			values += "?,"
			db.Args = append(db.Args, structValue)
		}

		values = values[:len(values)-1]
		db.SqlQuery += values + ")"
	}
}

// +----------------------------------------------------------------------
// | 更改方法
// +----------------------------------------------------------------------

func (db *PgQuery) UpdateByMap(data *map[string]interface{}) (string,[]interface{}) {
	db.OperateUpdateByMapData(data)
	db.replacePlace()
	return db.SqlQuery, db.Args
}
func (db *PgQuery) UpdateByStruct(data interface{}) (string,[]interface{}) {
	db.OperateUpdateByStructData(data)
	db.replacePlace()
	return db.SqlQuery, db.Args
}

//整理更改查询的sql和参数

func (db *PgQuery) OperateUpdateByMapData(data *map[string]interface{}) {
	numData := len(*data)
	if numData > 0 {
		db.SqlQuery += "UPDATE " + db.RecordTable + " "
		db.SqlQuery += "SET "
		var args []interface{}
		for key, val := range *data {
			db.SqlQuery += " " + key + " = ? ,"
			args = append(args, val)
		}
		db.Args = append(args, db.Args...)
		db.SqlQuery = db.SqlQuery[:len(db.SqlQuery)-1]
		if db.WhereSqlQuery != "" {
			db.SqlQuery += "where "
		}
		db.SqlQuery += db.WhereSqlQuery
	}
}
func (db *PgQuery) OperateUpdateByStructData(data interface{}) {
	dataType := reflect.TypeOf(data).Elem()
	dataValue := reflect.ValueOf(data).Elem()
	if dataType.Kind() != reflect.Struct {
		return
	}
	numField := dataValue.NumField()
	if numField > 0 {
		db.SqlQuery += "UPDATE " + db.RecordTable + " "
		db.SqlQuery += "SET "
		var args []interface{}
		for i := 0; i < dataValue.NumField(); i++ {
			field := dataType.Field(i).Tag.Get("json")
			structValue := dataValue.Field(i).Interface()
			db.SqlQuery += " " + field + " = ? ,"
			args = append(args, structValue)
		}
		db.SqlQuery = db.SqlQuery[:len(db.SqlQuery)-1]
		db.Args = append(args, db.Args...)
		if db.WhereSqlQuery != "" {
			db.SqlQuery += "where "
		}
		db.SqlQuery += db.WhereSqlQuery
	}
}

// +----------------------------------------------------------------------
// | 删除方法
// +----------------------------------------------------------------------

// Delete 删除方法
func (db *PgQuery) Delete() (string,[]interface{}) {
	db.OperateDeleteData()
	db.replacePlace()
	return db.SqlQuery, db.Args
}

// OperateDeleteData 整理删除查询的sql和参数
func (db *PgQuery) OperateDeleteData() {
	db.SqlQuery += "DELETE FROM  " + db.RecordTable + " "
	if db.WhereSqlQuery != "" {
		db.SqlQuery += "where "
	}
	db.SqlQuery += db.WhereSqlQuery
}

// +----------------------------------------------------------------------
// | 事务
// +----------------------------------------------------------------------

func PgTable(table string) *PgQuery {
	return &PgQuery{
		RecordTable: table,
	}
}

// +----------------------------------------------------------------------
// | my包没有的
// +----------------------------------------------------------------------

//因默认使用的?做占位符,pg数据库用的是$+变量,这里再生成sql语句的时候进行转换一下,变成pg可执行的sql
func (db *PgQuery) replacePlace(){
	split := strings.Split(db.SqlQuery, "?")
	db.SqlQuery = ""
	lenSplit := len(split)
	for i, s := range split {
		if i + 1 >= lenSplit{
			db.SqlQuery += fmt.Sprintf("%v ", s)
		}else{
			db.SqlQuery += fmt.Sprintf("%v $%v ", s,i+1)
		}
	}
}



