package sql_aid

import (
	"github.com/doujunyu/gogo/utility"
	"reflect"
	"strconv"
	"strings"
)

type MyQuery struct {
	RecordTable    string                 `json:"record_table" Testing:"表明"`
	RecordField    []string               `json:"record_field" Testing:"字段"`
	RecordOrder    []string               `json:"record_order" Testing:"排序"`
	RecordGroup    []string               `json:"record_group" Testing:"分组"`
	RecordPage     string                 `json:"record_page" Testing:"页数"`
	RecordSize     string                 `json:"record_size" Testing:"每页数据量"`
	RecordIncrease map[string]interface{} `json:"record_increase" Testing:"编辑中的增加数据"`
	RecordDecrease map[string]interface{} `json:"record_decrease" Testing:"编辑中的减少数据"`
	SqlQuery       string                 `json:"sql_query,string" Testing:"sql语句"`
	WhereSqlQuery  string                 `json:"where_sql_query" Testing:"sql条件"`
	Args           []interface{}          `json:"args" Testing:"值"`
}
type MyChildQuery func(*MyQuery, ...interface{})

// +----------------------------------------------------------------------
// | 查询
// +----------------------------------------------------------------------

//查询数据方法

func (db *MyQuery) ToSql() (string, []interface{}) {
	db.jointSql()
	return db.SqlQuery, db.Args
}
func (db *MyQuery) jointSql() (string, []interface{}) {
	db.OperateFindToSql()
	return db.SqlQuery, db.Args
}

//查询固定方法

func (db *MyQuery) Table(Table string) *MyQuery {
	db.RecordTable = "`" + Table + "`"
	return db
}
func (db *MyQuery) Field(field ...string) *MyQuery {
	for key, val := range field {
		field[key] = "`" + val + "`"
	}
	db.RecordField = field
	return db
}
func (db *MyQuery) Inc(field string,data interface{}) *MyQuery {
	db.RecordIncrease[field] = data
	return db
}
func (db *MyQuery) IncAll(incMap map[string]interface{}) *MyQuery {
	for key, val := range incMap {
		db.RecordIncrease[key] = val
	}
	return db
}
func (db *MyQuery) Dec(field string,data interface{}) *MyQuery {
	db.RecordDecrease[field] = data
	return db
}
func (db *MyQuery) DecAll(decMap map[string]interface{}) *MyQuery {
	for key, val := range decMap {
		db.RecordDecrease[key] = val
	}
	return db
}
func (db *MyQuery) OrderBy(Order string) *MyQuery {
	db.RecordOrder = append(db.RecordOrder, Order)
	return db
}
func (db *MyQuery) GroupBy(groupBy string) *MyQuery {
	db.RecordGroup = append(db.RecordGroup, groupBy)
	return db
}

//where条件

func (db *MyQuery) Where(field string, val interface{}) *MyQuery {
	if db.WhereSqlQuery != "" {
		db.WhereSqlQuery += "and "
	}
	db.WhereSqlQuery += field + " "
	if val != nil {
		db.Args = append(db.Args, val)
	}
	return db
}
func (db *MyQuery) WhereOr(field string, val interface{}) *MyQuery {
	if db.WhereSqlQuery = "" {
		db.WhereSqlQuery += "and "
	} else {
		db.WhereSqlQuery += "OR "
	}
	db.WhereSqlQuery += "(" + field + ") "
	if val != nil {
		db.Args = append(db.Args, val)
	}
	return db
}
func (db *MyQuery) WhereIn(field string, condition ...interface{}) *MyQuery {
	if db.WhereSqlQuery != "" {
		db.WhereSqlQuery += "and "
	}
	db.WhereSqlQuery += "`" + field + "` in ("
	for _, _ = range condition {
		db.WhereSqlQuery += "?,"
	}
	db.WhereSqlQuery = db.WhereSqlQuery[:len(db.WhereSqlQuery)-1]
	db.WhereSqlQuery += ") "
	db.Args = append(db.Args, condition...)
	return db
}
func (db *MyQuery) WhereNotIn(field string, condition ...interface{}) *MyQuery {
	if db.WhereSqlQuery != "" {
		db.WhereSqlQuery += "and "
	}
	db.WhereSqlQuery += "`" + field + "` not in ("
	for _, _ = range condition {
		db.WhereSqlQuery += "?,"
	}
	db.WhereSqlQuery = db.WhereSqlQuery[:len(db.WhereSqlQuery)-1]
	db.WhereSqlQuery += ") "
	db.Args = append(db.Args, condition...)
	return db
}
func (db *MyQuery) WhereRaw(childQuery MyChildQuery, val ...interface{}) *MyQuery {
	check := &MyQuery{}
	childQuery(check, val...)
	checkSql, args := check.jointSql()
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
func (db *MyQuery) WhereOrRaw(childQuery MyChildQuery, val ...interface{}) *MyQuery {
	check := &MyQuery{}
	childQuery(check, val...)
	checkSql, args := check.jointSql()
	if checkSql == "" {
		return db
	}
	if db.WhereSqlQuery != "" {
		db.WhereSqlQuery += "and ("
	} else {
		db.WhereSqlQuery += "( "
	}
	db.WhereSqlQuery += checkSql
	db.Args = append(db.Args, args...)
	db.WhereSqlQuery += ") "
	return db
}
func (db *MyQuery) WhereInRaw(field string, childQuery MyChildQuery, val ...interface{}) *MyQuery {
	if db.WhereSqlQuery != "" {
		db.WhereSqlQuery += "and "
	}
	db.WhereSqlQuery += "`" + field + "` in ("
	check := &MyQuery{}
	childQuery(check, val...)
	checkSql, args := check.jointSql()
	db.WhereSqlQuery += checkSql
	db.Args = append(db.Args, args...)
	db.WhereSqlQuery += ") "
	return db
}
func (db *MyQuery) WhereNotInRaw(field string, childQuery MyChildQuery, val ...interface{}) *MyQuery {
	if db.WhereSqlQuery != "" {
		db.WhereSqlQuery += "and "
	}
	db.WhereSqlQuery += "`" + field + "`not in ("
	check := &MyQuery{}
	childQuery(check, val...)
	checkSql, args := check.jointSql()
	db.WhereSqlQuery += checkSql
	db.Args = append(db.Args, args...)
	db.WhereSqlQuery += ") "
	return db
}
func (db *MyQuery) WhereId(id interface{}) *MyQuery {
	if db.WhereSqlQuery != "" {
		db.WhereSqlQuery += "and "
	}
	db.WhereSqlQuery = "`id` = ? "
	db.Args = append(db.Args, id)
	return db
}
func (db *MyQuery) PageSize(page string, size string) *MyQuery {
	if page == ""{
		page = "1"
	}
	if size == ""{
		size = "10"
	}
	db.RecordPage = page
	db.RecordSize = size
	return db
}

// 整理查询的sql和参数

func (db *MyQuery) OperateFindToSql() {
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
func (db *MyQuery) OperateFindField() {
	db.SqlQuery = "select "
	if db.RecordField != nil {
		db.SqlQuery += utility.StringBySliceString(",", db.RecordField) + " "
	} else {
		db.SqlQuery += "* "
	}
}
func (db *MyQuery) OperateFindTable() {
	db.SqlQuery += "FROM " + db.RecordTable + " "
}
func (db *MyQuery) OperateFindGroupBy() {
	if db.RecordGroup != nil {
		db.SqlQuery += "GROUP BY " + strings.Join(db.RecordGroup, ",") + " "
	}
}
func (db *MyQuery) OperateFindOrderBy() {
	if db.RecordOrder != nil {
		db.SqlQuery += "ORDER BY " + strings.Join(db.RecordOrder, ",") + " "
	}
}

func (db *MyQuery) OperateFindPageSize() {
	page, err := strconv.Atoi(db.RecordPage)
	if err != nil {
		return
	}
	if page <= 0 {
		page = 1
	}
	size, _ := strconv.Atoi(db.RecordSize)
	if size == 0 {
		size = 10
	}
	db.SqlQuery += "limit ?,? "
	db.Args = append(db.Args, (page-1)*size)
	db.Args = append(db.Args, size)

}

// +----------------------------------------------------------------------
// | 添加方法
// +----------------------------------------------------------------------

func (db *MyQuery) InsertByMap(data map[string]interface{}) (string, []interface{}) {
	db.OperateInsertTable()
	db.OperateInsertDataByMap(data)
	return db.SqlQuery, db.Args
}
func (db *MyQuery) InsertByStruct(data interface{}) (string, []interface{}) {
	db.OperateInsertTable()
	db.OperateInsertDataByStruct(data)
	return db.SqlQuery, db.Args
}
func (db *MyQuery) InsertAllByMap(datas *[]map[string]interface{}) (string, []interface{}) {
	db.OperateInsertTable()
	dataSlice := *datas
	if len(dataSlice) < 1{
		return "",nil
	}
	fieldSlice := db.OperateInsertDataByMap(dataSlice[0])
	for _, val := range dataSlice[1:] {
		db.OperateInsertDataByMapValue(val,fieldSlice)
	}

	return db.SqlQuery, db.Args

}
func (db *MyQuery) InsertAllByStruct(datas []interface{}) (string, []interface{}) {
	db.OperateInsertTable()
	for key, val := range datas {
		if key == 0 {
			db.OperateInsertDataByStruct(val)
		} else {
			db.OperateInsertDataByStructValue(val)
		}
	}
	return db.SqlQuery, db.Args

}

// 整理查询的sql和参数

func (db *MyQuery) OperateInsertTable() {
	if db.RecordTable != "" {
		db.SqlQuery += "INSERT INTO `" + db.RecordTable + "` "
	}
}
func (db *MyQuery) OperateInsertDataByMap(data map[string]interface{}) *[]string {
	numData := len(data)
	field := make([]string,0)
	if numData > 0 {
		db.SqlQuery += "("
		values := ""
		for key, val := range data {
			field = append(field,key)
			db.SqlQuery += "`" + key + "`,"
			values += "?,"
			db.Args = append(db.Args, val)
		}
		db.SqlQuery = db.SqlQuery[:len(db.SqlQuery)-1]
		values = values[:len(values)-1]
		db.SqlQuery += ") VALUES (" + values + ")"
	}
	return &field
}
func (db *MyQuery) OperateInsertDataByMapValue(data map[string]interface{},fieldSlice *[]string) {
	db.SqlQuery += ",( "
	values := ""
	for _, val := range *fieldSlice {
		values += "?,"
		db.Args = append(db.Args, data[val])
	}
	values = values[:len(values)-1]
	db.SqlQuery += values + ")"
}
func (db *MyQuery) OperateInsertDataByStruct(data interface{}) {
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
			db.SqlQuery += "`" + field + "`,"
			values += "?,"
			db.Args = append(db.Args, structValue)
		}
		db.SqlQuery = db.SqlQuery[:len(db.SqlQuery)-1]
		values = values[:len(values)-1]
		db.SqlQuery += ") VALUES (" + values + ") "
	}
}
func (db *MyQuery) OperateInsertDataByStructValue(data interface{}) {
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

func (db *MyQuery) UpdateByMap(data map[string]interface{}) (string, []interface{}) {
	db.OperateUpdateByMapData(data)
	return db.SqlQuery, db.Args
}
func (db *MyQuery) UpdateByStruct(data interface{}) (string, []interface{}) {
	db.OperateUpdateByStructData(data)
	return db.SqlQuery, db.Args
}

//整理更改查询的sql和参数

func (db *MyQuery) OperateUpdateByMapData(data map[string]interface{}) {
	db.SqlQuery += "UPDATE `" + db.RecordTable + "` "
	db.SqlQuery += "SET "
	var args []interface{}
	for key, val := range data {
		db.SqlQuery += " `" + key + "` = ? ,"
		args = append(args, val)
	}
	for incKey, incVal := range db.RecordIncrease {
		db.SqlQuery += " " + incKey + " = "+incKey+" + ? ,"
		args = append(args, incVal)
	}
	for decKey, decVal := range db.RecordDecrease {
		db.SqlQuery += " " + decKey + " = "+decKey+" - ? ,"
		args = append(args, decVal)
	}
	db.Args = append(args, db.Args...)
	db.SqlQuery = db.SqlQuery[:len(db.SqlQuery)-1]
	if db.WhereSqlQuery != "" {
		db.SqlQuery += "where "
	}
	db.SqlQuery += db.WhereSqlQuery
}
func (db *MyQuery) OperateUpdateByStructData(data interface{}) {
	dataType := reflect.TypeOf(data).Elem()
	dataValue := reflect.ValueOf(data).Elem()
	if dataType.Kind() != reflect.Struct {
		return
	}
	db.SqlQuery += "UPDATE `" + db.RecordTable + "` "
	db.SqlQuery += "SET "
	var args []interface{}
	for i := 0; i < dataValue.NumField(); i++ {
		field := dataType.Field(i).Tag.Get("json")
		structValue := dataValue.Field(i).Interface()
		db.SqlQuery += " `" + field + "` = ? ,"
		args = append(args, structValue)
	}
	for incKey, incVal := range db.RecordIncrease {
		db.SqlQuery += " " + incKey + " = "+incKey+" + ? ,"
		args = append(args, incVal)
	}
	for decKey, decVal := range db.RecordDecrease {
		db.SqlQuery += " " + decKey + " = "+decKey+" - ? ,"
		args = append(args, decVal)
	}
	db.SqlQuery = db.SqlQuery[:len(db.SqlQuery)-1]
	db.Args = append(args, db.Args...)
	if db.WhereSqlQuery != "" {
		db.SqlQuery += "where "
	}
	db.SqlQuery += db.WhereSqlQuery
}

// +----------------------------------------------------------------------
// | 删除方法
// +----------------------------------------------------------------------

// Delete 删除方法
func (db *MyQuery) Delete() (string, []interface{}) {
	db.OperateDeleteData()
	return db.SqlQuery, db.Args
}

// OperateDeleteData 整理删除查询的sql和参数
func (db *MyQuery) OperateDeleteData() {
	db.SqlQuery += "DELETE FROM  `" + db.RecordTable + "` "
	if db.WhereSqlQuery != "" {
		db.SqlQuery += "where "
	}
	db.SqlQuery += db.WhereSqlQuery
}

// +----------------------------------------------------------------------
// | 事务
// +----------------------------------------------------------------------

func MyTable(table string) *MyQuery {
	return &MyQuery{
		RecordTable: table,
		RecordIncrease: make(map[string]interface{}),
		RecordDecrease: make(map[string]interface{}),
	}
}
