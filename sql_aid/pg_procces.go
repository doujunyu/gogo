package sql_aid

import (
	"fmt"
	"github.com/doujunyu/gogo/utility"
	"reflect"
	"strconv"
	"strings"
)

type PgQuery struct {
	recordTable    string                 `Testing:"表明"`
	recordField    []string               `Testing:"字段"`
	recordOrder    []string               `Testing:"排序"`
	recordGroup    []string               `Testing:"分组"`
	recordPage     string                 `Testing:"页数"`
	recordSize     string                 `Testing:"每页数据量"`
	recordIncrease map[string]interface{} `Testing:"编辑中的增加数据"`
	recordDecrease map[string]interface{} `Testing:"编辑中的减少数据"`
	sqlQuery       string                 `Testing:"sql语句"`
	whereSqlQuery  string                 `Testing:"sql条件"`
	args           []interface{}          `Testing:"值"`
}
type PgChildQuery func(*PgQuery, ...interface{})

// +----------------------------------------------------------------------
// | 查询
// +----------------------------------------------------------------------

//查询数据方法

func (db *PgQuery) ToSql() (string, []interface{}) {
	db.jointSql()
	db.replacePlace()
	return db.sqlQuery, db.args
}
func (db *PgQuery) jointSql() (string, []interface{}) {
	db.operateFindToSql()
	return db.sqlQuery, db.args
}

//查询固定方法

func (db *PgQuery) Table(Table string) *PgQuery {
	db.recordTable = Table
	return db
}
func (db *PgQuery) Field(field ...string) *PgQuery {
	for key, val := range field {
		field[key] = val
	}
	db.recordField = field
	return db
}
func (db *PgQuery) Inc(field string,data interface{}) *PgQuery {
	db.recordIncrease[field] = data
	return db
}
func (db *PgQuery) IncAll(incMap map[string]interface{}) *PgQuery {
	for key, val := range incMap {
		db.recordIncrease[key] = val
	}
	return db
}
func (db *PgQuery) Dec(field string,data interface{}) *PgQuery {
	db.recordDecrease[field] = data
	return db
}

func (db *PgQuery) DecAll(decMap map[string]interface{}) *PgQuery {
	for key, val := range decMap {
		db.recordDecrease[key] = val
	}
	return db
}

func (db *PgQuery) OrderBy(Order string) *PgQuery {
	db.recordOrder = append(db.recordOrder, Order)
	return db
}
func (db *PgQuery) GroupBy(groupBy string) *PgQuery {
	db.recordGroup = append(db.recordGroup, groupBy)
	return db
}

//where条件

func (db *PgQuery) Where(field string, val interface{}) *PgQuery {
	if db.whereSqlQuery != "" {
		db.whereSqlQuery += "and "
	}
	db.whereSqlQuery += field + " "
	if val != nil {
		db.args = append(db.args, val)
	}
	return db
}
func (db *PgQuery) WhereOr(field string, val interface{}) *PgQuery {
	if db.whereSqlQuery == "" {
		db.whereSqlQuery += "and "
	} else {
		db.whereSqlQuery += "OR "
	}
	db.whereSqlQuery += "(" + field + ") "
	if val != nil {
		db.args = append(db.args, val)
	}
	return db
}
func (db *PgQuery) WhereIn(field string, condition ...interface{}) *PgQuery {
	if db.whereSqlQuery != "" {
		db.whereSqlQuery += "and "
	}
	db.whereSqlQuery += field + " in ("
	for _, _ = range condition {
		db.whereSqlQuery += "?,"
	}
	db.whereSqlQuery = db.whereSqlQuery[:len(db.whereSqlQuery)-1]
	db.whereSqlQuery += ") "
	db.args = append(db.args, condition...)

	return db
}
func (db *PgQuery) WhereNotIn(field string, condition ...interface{}) *PgQuery {
	if db.whereSqlQuery != "" {
		db.whereSqlQuery += "and "
	}
	db.whereSqlQuery += field + " not in ("
	for _, _ = range condition {
		db.whereSqlQuery += "?,"
	}
	db.whereSqlQuery = db.whereSqlQuery[:len(db.whereSqlQuery)-1]
	db.whereSqlQuery += ") "
	db.args = append(db.args, condition...)

	return db
}
func (db *PgQuery) WhereRaw(childQuery PgChildQuery, val ...interface{}) *PgQuery {
	check := &PgQuery{}
	childQuery(check, val...)
	checkSql, args := check.jointSql()
	if checkSql == "" {
		return db
	}
	if db.whereSqlQuery != "" {
		db.whereSqlQuery += "and "
	}
	db.whereSqlQuery += "("
	db.whereSqlQuery += checkSql
	db.args = append(db.args, args...)
	db.whereSqlQuery += ") "
	return db
}
func (db *PgQuery) WhereOrRaw(childQuery PgChildQuery, val ...interface{}) *PgQuery {
	check := &PgQuery{}
	childQuery(check, val...)
	checkSql, args := check.jointSql()
	if checkSql == "" {
		return db
	}
	if db.whereSqlQuery != "" {
		db.whereSqlQuery += "OR ("
	} else {
		db.whereSqlQuery += "( "
	}
	db.whereSqlQuery += checkSql
	db.args = append(db.args, args...)
	db.whereSqlQuery += ") "
	return db
}
func (db *PgQuery) WhereInRaw(field string, childQuery PgChildQuery, val ...interface{}) *PgQuery {
	if db.whereSqlQuery != "" {
		db.whereSqlQuery += "and "
	}
	db.whereSqlQuery += field + " in ("
	check := &PgQuery{}
	childQuery(check, val...)
	checkSql, args := check.jointSql()
	db.whereSqlQuery += checkSql
	db.args = append(db.args, args...)
	db.whereSqlQuery += ") "
	return db
}
func (db *PgQuery) WhereNotInRaw(field string, childQuery PgChildQuery, val ...interface{}) *PgQuery {
	if db.whereSqlQuery != "" {
		db.whereSqlQuery += "and "
	}
	db.whereSqlQuery += field + "not in ("
	check := &PgQuery{}
	childQuery(check, val...)
	checkSql, args := check.jointSql()
	db.whereSqlQuery += checkSql
	db.args = append(db.args, args...)
	db.whereSqlQuery += ") "
	return db
}
func (db *PgQuery) WhereId(id interface{}) *PgQuery {
	if db.whereSqlQuery != "" {
		db.whereSqlQuery += "and "
	}
	db.whereSqlQuery = "id = ? "
	db.args = append(db.args, id)
	return db
}
func (db *PgQuery) PageSize(page string, size string) *PgQuery {
	if page == ""{
		page = "1"
	}
	if size == ""{
		size = "10"
	}
	db.recordPage = page
	db.recordSize = size
	return db
}

// 整理查询的sql和参数

func (db *PgQuery) operateFindToSql() {
	if db.recordTable != "" {
		db.operateFindField()
		db.operateFindTable()
	}
	if db.recordTable != "" && db.whereSqlQuery != "" {
		db.sqlQuery += "where "
	}
	db.sqlQuery += db.whereSqlQuery
	db.operateFindGroupBy()
	db.operateFindOrderBy()
	db.operateFindPageSize()
}
func (db *PgQuery) operateFindField() {
	db.sqlQuery = "select "
	if db.recordField != nil {
		db.sqlQuery += utility.StringBySliceString(",", db.recordField) + " "
	} else {
		db.sqlQuery += "* "
	}
}
func (db *PgQuery) operateFindTable() {
	db.sqlQuery += "FROM " + db.recordTable + " "
}
func (db *PgQuery) operateFindGroupBy() {
	if db.recordGroup != nil {
		db.sqlQuery += "GROUP BY " + strings.Join(db.recordGroup, ",") + " "
	}
}
func (db *PgQuery) operateFindOrderBy() {
	if db.recordOrder != nil {
		db.sqlQuery += "ORDER BY " + strings.Join(db.recordOrder, ",") + " "
	}
}

func (db *PgQuery) operateFindPageSize() {
	page, err := strconv.Atoi(db.recordPage)
	if err != nil {
		return
	}
	if page <= 0 {
		page = 1
	}
	size, _ := strconv.Atoi(db.recordSize)
	if size == 0 {
		size = 10
	}
	db.sqlQuery += "limit ? OFFSET ? "
	db.args = append(db.args, size)
	db.args = append(db.args, (page-1)*size)

}

// +----------------------------------------------------------------------
// | 添加方法
// +----------------------------------------------------------------------

func (db *PgQuery) InsertByMap(data map[string]interface{}) (string, []interface{}) {
	db.operateInsertTable()
	db.operateInsertDataByMap(data)
	db.replacePlace()
	return db.sqlQuery, db.args
}
func (db *PgQuery) InsertByStruct(data interface{}) (string, []interface{}) {
	db.operateInsertTable()
	db.operateInsertDataByStruct(data)
	db.replacePlace()
	return db.sqlQuery, db.args
}
func (db *PgQuery) InsertAllByMap(datas *[]map[string]interface{}) (string, []interface{}) {
	db.operateInsertTable()
	dataSlice := *datas
	if len(dataSlice) < 1{
		return "",nil
	}
	fieldSlice := db.operateInsertDataByMap(dataSlice[0])
	for _, val := range dataSlice[1:] {
		db.operateInsertDataByMapValue(val,fieldSlice)
	}
	db.replacePlace()
	return db.sqlQuery, db.args

}
func (db *PgQuery) InsertAllByStruct(datas []interface{}) (string, []interface{}) {
	db.operateInsertTable()
	for key, val := range datas {
		if key == 0 {
			db.operateInsertDataByStruct(val)
		} else {
			db.operateInsertDataByStructValue(val)
		}
	}
	db.replacePlace()
	return db.sqlQuery, db.args

}

// 整理查询的sql和参数

func (db *PgQuery) operateInsertTable() {
	if db.recordTable != "" {
		db.sqlQuery += "INSERT INTO " + db.recordTable + " "
	}
}
func (db *PgQuery) operateInsertDataByMap(data map[string]interface{}) *[]string {
	numData := len(data)
	field := make([]string,0)
	if numData > 0 {
		db.sqlQuery += "("
		values := ""
		for key, val := range data {
			field = append(field,key)
			db.sqlQuery += key + ","
			values += "?,"
			db.args = append(db.args, val)
		}
		db.sqlQuery = db.sqlQuery[:len(db.sqlQuery)-1]
		values = values[:len(values)-1]
		db.sqlQuery += ") VALUES (" + values + ")"
	}
	return &field
}
func (db *PgQuery) operateInsertDataByMapValue(data map[string]interface{},fieldSlice *[]string) {
	db.sqlQuery += ",( "
	values := ""
	for _, val := range *fieldSlice {
		values += "?,"
		db.args = append(db.args, data[val])
	}
	values = values[:len(values)-1]
	db.sqlQuery += values + ")"
}
func (db *PgQuery) operateInsertDataByStruct(data interface{}) {
	dataType := reflect.TypeOf(data).Elem()
	dataValue := reflect.ValueOf(data).Elem()
	if dataType.Kind() != reflect.Struct {
		return
	}
	numField := dataValue.NumField()
	if numField > 0 {
		values := ""
		db.sqlQuery += "("
		for i := 0; i < dataValue.NumField(); i++ {
			field := dataType.Field(i).Tag.Get("json")
			structValue := dataValue.Field(i).Interface()
			db.sqlQuery += field + ","
			values += "?,"
			db.args = append(db.args, structValue)
		}
		db.sqlQuery = db.sqlQuery[:len(db.sqlQuery)-1]
		values = values[:len(values)-1]
		db.sqlQuery += ") VALUES (" + values + ") "
	}
}
func (db *PgQuery) operateInsertDataByStructValue(data interface{}) {
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
			db.args = append(db.args, structValue)
		}

		values = values[:len(values)-1]
		db.sqlQuery += values + ")"
	}
}

// +----------------------------------------------------------------------
// | 更改方法
// +----------------------------------------------------------------------

func (db *PgQuery) UpdateByMap(data map[string]interface{}) (string, []interface{}) {
	db.operateUpdateByMapData(data)
	db.replacePlace()
	return db.sqlQuery, db.args
}
func (db *PgQuery) UpdateByStruct(data interface{}) (string, []interface{}) {
	db.operateUpdateByStructData(data)
	db.replacePlace()
	return db.sqlQuery, db.args
}

//整理更改查询的sql和参数

func (db *PgQuery) operateUpdateByMapData(data map[string]interface{}) {
	db.sqlQuery += "UPDATE " + db.recordTable + " "
	db.sqlQuery += "SET "
	var args []interface{}
	for key, val := range data {
		db.sqlQuery += " " + key + " = ? ,"
		args = append(args, val)
	}
	for incKey, incVal := range db.recordIncrease {
		db.sqlQuery += " " + incKey + " = "+incKey+" + ? ,"
		args = append(args, incVal)
	}
	for decKey, decVal := range db.recordDecrease {
		db.sqlQuery += " " + decKey + " = "+decKey+" - ? ,"
		args = append(args, decVal)
	}
	db.args = append(args, db.args...)
	db.sqlQuery = db.sqlQuery[:len(db.sqlQuery)-1]
	if db.whereSqlQuery != "" {
		db.sqlQuery += "where "
	}
	db.sqlQuery += db.whereSqlQuery
}
func (db *PgQuery) operateUpdateByStructData(data interface{}) {
	dataType := reflect.TypeOf(data).Elem()
	dataValue := reflect.ValueOf(data).Elem()
	if dataType.Kind() != reflect.Struct {
		return
	}
		db.sqlQuery += "UPDATE " + db.recordTable + " "
		db.sqlQuery += "SET "
		var args []interface{}
		for i := 0; i < dataValue.NumField(); i++ {
			field := dataType.Field(i).Tag.Get("json")
			structValue := dataValue.Field(i).Interface()
			db.sqlQuery += " " + field + " = ? ,"
			args = append(args, structValue)
		}
		for incKey, incVal := range db.recordIncrease {
			db.sqlQuery += " " + incKey + " = "+incKey+" + ? ,"
			args = append(args, incVal)
		}
		for decKey, decVal := range db.recordDecrease {
			db.sqlQuery += " " + decKey + " = "+decKey+" - ? ,"
			args = append(args, decVal)
		}
		db.sqlQuery = db.sqlQuery[:len(db.sqlQuery)-1]
		db.args = append(args, db.args...)
		if db.whereSqlQuery != "" {
			db.sqlQuery += "where "
		}
		db.sqlQuery += db.whereSqlQuery

}

// +----------------------------------------------------------------------
// | 删除方法
// +----------------------------------------------------------------------

// Delete 删除方法
func (db *PgQuery) Delete() (string, []interface{}) {
	db.operateDeleteData()
	db.replacePlace()
	return db.sqlQuery, db.args
}

// OperateDeleteData 整理删除查询的sql和参数
func (db *PgQuery) operateDeleteData() {
	db.sqlQuery += "DELETE FROM  " + db.recordTable + " "
	if db.whereSqlQuery != "" {
		db.sqlQuery += "where "
	}
	db.sqlQuery += db.whereSqlQuery
}

// +----------------------------------------------------------------------
// | 开头
// +----------------------------------------------------------------------

func PgTable(table string) *PgQuery {
	return &PgQuery{
		recordTable: table,
		recordIncrease: make(map[string]interface{}),
		recordDecrease: make(map[string]interface{}),
	}
}

// +----------------------------------------------------------------------
// | my包没有的
// +----------------------------------------------------------------------

//因默认使用的?做占位符,pg数据库用的是$+变量,这里再生成sql语句的时候进行转换一下,变成pg可执行的sql
func (db *PgQuery) replacePlace() {
	split := strings.Split(db.sqlQuery, "?")
	db.sqlQuery = ""
	lenSplit := len(split)
	for i, s := range split {
		if i+1 >= lenSplit {
			db.sqlQuery += fmt.Sprintf("%v", s)
		} else {
			db.sqlQuery += fmt.Sprintf("%v$%v", s, i+1)
		}
	}
}
