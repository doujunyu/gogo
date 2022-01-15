package sql_aid

import (
	"github.com/doujunyu/gogo/utility"
	"reflect"
	"strings"
)

type Query struct {
	RecordTable   string        `json:"record_table" Testing:"表明"`
	RecordField   []string      `json:"record_field" Testing:"字段"`
	RecordOrder   []string        `json:"record_order" Testing:"排序"`
	RecordGroup   []string        `json:"record_group" Testing:"分组"`
	RecordPage    int           `json:"record_page" Testing:"页数"`
	RecordSize    int           `json:"record_size" Testing:"每页数据量"`
	SqlQuery      string        `json:"sql_query,string" Testing:"sql语句"`
	WhereSqlQuery string        `json:"where_sql_query" Testing:"sql条件"`
	Args          []interface{} `json:"args" Testing:"值"`

}
type ChildQuery func(*Query, ...interface{})

// +----------------------------------------------------------------------
// | 查询
// +----------------------------------------------------------------------

//查询数据方法

func (db *Query) ToSql() (string, []interface{}) {
	db.OperateFindToSql()
	return db.SqlQuery, db.Args
}

//查询固定方法

func (db *Query) Table(Table string) *Query {
	db.RecordTable = "`" + Table + "`"
	return db
}
func (db *Query) Field(field ...string) *Query {
	for key, val := range field {
		field[key] = "`" + val + "`"
	}
	db.RecordField = field
	return db
}
func (db *Query) OrderBy(Order string) *Query {
	db.RecordOrder = append(db.RecordOrder,Order)
	return db
}
func (db *Query) GroupBy(groupBy string) *Query {
	db.RecordGroup = append(db.RecordGroup,groupBy)
	return db
}

//where条件

func (db *Query) Where(field string, val interface{}) *Query {
	if db.WhereSqlQuery != "" {
		db.WhereSqlQuery += "and "
	}
	db.WhereSqlQuery += field
	db.Args = append(db.Args, val)
	return db
}
func (db *Query) WhereOr(field string, val interface{}) *Query {
	if db.WhereSqlQuery != "" {
		db.WhereSqlQuery += "and "
	}else{
		db.WhereSqlQuery += "OR "
	}
	db.WhereSqlQuery += "(" + field + ") "
	db.Args = append(db.Args, val)
	return db
}
func (db *Query) WhereIn(field string, condition ...interface{}) *Query {
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
func (db *Query) WhereNotIn(field string, condition ...interface{}) *Query {
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
func (db *Query) WhereRaw(childQuery ChildQuery, val ...interface{}) *Query {
	if db.WhereSqlQuery != "" {
		db.WhereSqlQuery += "and "
	}
	db.WhereSqlQuery += "("
	check := &Query{}
	childQuery(check, val...)
	checkSql,args := check.ToSql()
	db.WhereSqlQuery += checkSql
	db.Args = append(db.Args, args...)
	db.WhereSqlQuery += ") "
	return db
}
func (db *Query) WhereOrRaw(childQuery ChildQuery, val ...interface{}) *Query {
	if db.WhereSqlQuery != "" {
		db.WhereSqlQuery += "and ("
	}else{
		db.WhereSqlQuery += "( "
	}
	check := &Query{}
	childQuery(check, val...)
	checkSql,args := check.ToSql()
	db.WhereSqlQuery += checkSql
	db.Args = append(db.Args, args...)
	db.WhereSqlQuery += ") "
	return db
}
func (db *Query) WhereInRaw(field string, childQuery ChildQuery, val ...interface{}) *Query {
	if db.WhereSqlQuery != "" {
		db.WhereSqlQuery += "and "
	}
	db.WhereSqlQuery += "`" + field + "` in ("
	check := &Query{}
	childQuery(check, val...)
	checkSql,args := check.ToSql()
	db.WhereSqlQuery += checkSql
	db.Args = append(db.Args, args...)
	db.WhereSqlQuery += ") "
	return db
}
func (db *Query) WhereNotInRaw(field string, childQuery ChildQuery, val ...interface{}) *Query {
	if db.WhereSqlQuery != "" {
		db.WhereSqlQuery += "and "
	}
	db.WhereSqlQuery += "`" + field + "`not in ("
	check := &Query{}
	childQuery(check, val...)
	checkSql,args := check.ToSql()
	db.WhereSqlQuery += checkSql
	db.Args = append(db.Args, args...)
	db.WhereSqlQuery += ") "
	return db
}
func (db *Query) WhereId(id string) *Query {
	if db.WhereSqlQuery != "" {
		db.WhereSqlQuery += "and "
	}
	db.WhereSqlQuery = "`id` = ? "
	db.Args = append(db.Args, id)
	return db
}
func (db *Query) PageSize(page int, size int) *Query {
	db.RecordPage = page
	db.RecordSize = size
	return db
}

// 整理查询的sql和参数

func (db *Query) OperateFindToSql() {
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
func (db *Query) OperateFindField() {
	db.SqlQuery = "select "
	if db.RecordField != nil {
		db.SqlQuery += utility.StringBySliceString(",", db.RecordField) + " "
	} else {
		db.SqlQuery += "* "
	}
}
func (db *Query) OperateFindTable() {
	db.SqlQuery += "FROM " + db.RecordTable + " "
}
func (db *Query) OperateFindGroupBy() {
	if db.RecordGroup != nil{
		db.SqlQuery += "GROUP BY " + strings.Join(db.RecordGroup,",") + " "
	}
}
func (db *Query) OperateFindOrderBy() {
	if db.RecordOrder != nil{
		db.SqlQuery += "ORDER BY " + strings.Join(db.RecordOrder,",") + " "
	}
}

func (db *Query) OperateFindPageSize() {
	if db.RecordPage != 0 {
		if db.RecordSize == 0 {
			db.RecordSize = 10
		}
		var limita int = (db.RecordPage - 1) * db.RecordSize
		db.SqlQuery += "limit ?,? "
		db.Args = append(db.Args, limita)
		db.Args = append(db.Args, db.RecordSize)
	}

}

// +----------------------------------------------------------------------
// | 添加方法
// +----------------------------------------------------------------------

func (db *Query) InsertByMap(data *map[string]interface{}) (string,[]interface{}) {
	db.OperateInsertTable()
	db.OperateInsertDataByMap(data)
	return db.SqlQuery, db.Args
}
func (db *Query) InsertByStruct(data interface{}) (string,[]interface{}) {
	db.OperateInsertTable()
	db.OperateInsertDataByStruct(data)
	return db.SqlQuery, db.Args
}
func (db *Query) InsertAllByMap(datas *[]map[string]interface{}) (string,[]interface{}) {
	db.OperateInsertTable()
	for key, val := range *datas {
		if key == 0 {
			db.OperateInsertDataByMap(&val)
		} else {
			db.OperateInsertDataByMapValue(&val)
		}
	}
	return db.SqlQuery, db.Args

}
func (db *Query) InsertAllByStruct(datas []interface{}) (string,[]interface{}) {
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

func (db *Query) OperateInsertTable() {
	if db.RecordTable != "" {
		db.SqlQuery += "INSERT INTO `" + db.RecordTable + "` "
	}
}
func (db *Query) OperateInsertDataByMap(data *map[string]interface{}) {
	numData := len(*data)
	if numData > 0 {
		db.SqlQuery += "("
		values := ""
		for key, val := range *data {
			db.SqlQuery += "`" + key + "`,"
			values += "?,"
			db.Args = append(db.Args, val)
		}
		db.SqlQuery = db.SqlQuery[:len(db.SqlQuery)-1]
		values = values[:len(values)-1]
		db.SqlQuery += ") VALUES (" + values + ")"
	}

}
func (db *Query) OperateInsertDataByMapValue(data *map[string]interface{}) {
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
func (db *Query) OperateInsertDataByStruct(data interface{}) {
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
func (db *Query) OperateInsertDataByStructValue(data interface{}) {
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

func (db *Query) UpdateByMap(data *map[string]interface{}) (string,[]interface{}) {
	db.OperateUpdateByMapData(data)
	return db.SqlQuery, db.Args
}
func (db *Query) UpdateByStruct(data interface{}) (string,[]interface{}) {
	db.OperateUpdateByStructData(data)
	return db.SqlQuery, db.Args
}

//整理更改查询的sql和参数

func (db *Query) OperateUpdateByMapData(data *map[string]interface{}) {
	numData := len(*data)
	if numData > 0 {
		db.SqlQuery += "UPDATE `" + db.RecordTable + "` "
		db.SqlQuery += "SET "
		var args []interface{}
		for key, val := range *data {
			db.SqlQuery += " `" + key + "` = ? ,"
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
func (db *Query) OperateUpdateByStructData(data interface{}) {
	dataType := reflect.TypeOf(data).Elem()
	dataValue := reflect.ValueOf(data).Elem()
	if dataType.Kind() != reflect.Struct {
		return
	}
	numField := dataValue.NumField()
	if numField > 0 {
		db.SqlQuery += "UPDATE `" + db.RecordTable + "` "
		db.SqlQuery += "SET "
		var args []interface{}
		for i := 0; i < dataValue.NumField(); i++ {
			field := dataType.Field(i).Tag.Get("json")
			structValue := dataValue.Field(i).Interface()
			db.SqlQuery += " `" + field + "` = ? ,"
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
func (db *Query) Delete() (string,[]interface{}) {
	db.OperateDeleteData()
	return db.SqlQuery, db.Args
}

// OperateDeleteData 整理删除查询的sql和参数
func (db *Query) OperateDeleteData() {
	db.SqlQuery += "DELETE FROM  `" + db.RecordTable + "` "
	if db.WhereSqlQuery != "" {
		db.SqlQuery += "where "
	}
	db.SqlQuery += db.WhereSqlQuery
}

// +----------------------------------------------------------------------
// | 事务
// +----------------------------------------------------------------------

func Table(table string) *Query {
	return &Query{
		RecordTable: table,
	}
}