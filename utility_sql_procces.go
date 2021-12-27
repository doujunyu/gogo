package gogo

import (
	"database/sql"
	"fmt"
	"reflect"
)

type Query struct {
	RecordTable   string        `json:"record_table" Testing:"表明"`
	RecordField   []string      `json:"record_field" Testing:"字段"`
	RecordOrder   string        `json:"record_order" Testing:"排序"`
	RecordGroup   string        `json:"record_group" Testing:"分组"`
	RecordPage    int           `json:"record_page" Testing:"页数"`
	RecordSize    int           `json:"record_size" Testing:"每页数据量"`
	SqlQuery      string        `json:"sql_query,string" Testing:"sql语句"`
	WhereSqlQuery string        `json:"where_sql_query" Testing:"sql条件"`
	Args          []interface{} `json:"args" Testing:"值"`
	Tx            *sql.Tx
}
type ChildQuery func(*Query, ...interface{})

// +----------------------------------------------------------------------
// | 查询
// +----------------------------------------------------------------------

//查询数据方法

func (db *Query) Find() ([]map[string]interface{}, error) {
	db.OperateFindToSql()
	var rows *sql.Rows
	var err error
	if db.Tx != nil {
		rows, err = db.Tx.Query(db.SqlQuery, db.Args...)
	} else {
		rows, err = Open().Query(db.SqlQuery, db.Args...)
	}
	if err != nil {
		return nil, err
	}
	return QueryFind(rows)
}
func (db *Query) FindOnly() (map[string]interface{}, error) {
	data, err := db.Find()
	if err != nil {
		return nil, err
	}
	return data[0], nil
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
	db.RecordOrder = "order by " + Order + " "
	return db
}
func (db *Query) GroupBy(groupBy string) *Query {
	db.RecordOrder = "GROUP BY `" + groupBy + "` "
	return db
}

//where条件

func (db *Query) Where(field string, condition string, val interface{}) *Query {
	if db.WhereSqlQuery != "" {
		db.WhereSqlQuery += "and "
	}
	db.WhereSqlQuery += "`" + field + "` " + condition + " ? "
	db.Args = append(db.Args, val)
	return db
}
func (db *Query) WhereRaw(childQuery ChildQuery, val ...interface{}) *Query {
	db.WhereSqlQuery += "and ("
	check := &Query{}
	childQuery(check, val...)
	if check.SqlQuery != "" {
		check.SqlQuery += check.SqlQuery + "where " + check.WhereSqlQuery
	} else {
		check.SqlQuery += check.WhereSqlQuery
	}
	db.WhereSqlQuery += check.SqlQuery
	for _, val := range check.Args {
		db.Args = append(db.Args, val)
	}
	db.WhereSqlQuery += ") "
	return db
}
func (db *Query) WhereOr(field string, condition string, val interface{}) *Query {
	db.WhereSqlQuery += "OR (`" + field + "` " + condition + " ?) "
	db.Args = append(db.Args, val)
	return db
}
func (db *Query) WhereOrRaw(childQuery ChildQuery, val ...interface{}) *Query {
	db.WhereSqlQuery += "OR ("
	check := &Query{}
	childQuery(check, val...)
	check.OperateFindToSql()
	db.WhereSqlQuery += check.SqlQuery
	for _, val := range check.Args {
		db.Args = append(db.Args, val)
	}
	db.WhereSqlQuery += ") "
	return db
}
func (db *Query) WhereIn(field string, condition ...interface{}) *Query {
	if db.WhereSqlQuery != "" {
		db.WhereSqlQuery += "and "
	}
	db.WhereSqlQuery += "`" + field + "` in ("
	for _, val := range condition {
		db.WhereSqlQuery += "?,"
		db.Args = append(db.Args, val)
	}
	db.WhereSqlQuery = db.WhereSqlQuery[:len(db.WhereSqlQuery)-1]
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
	check.OperateFindToSql()
	db.WhereSqlQuery += check.SqlQuery
	for _, val := range check.Args {
		db.Args = append(db.Args, val)
	}
	db.WhereSqlQuery += ") "
	return db
}
func (db *Query) WhereNotIn(field string, condition ...interface{}) *Query {
	if db.WhereSqlQuery != "" {
		db.WhereSqlQuery += "and "
	}
	db.WhereSqlQuery += "`" + field + "` not in ("
	for _, val := range condition {
		db.WhereSqlQuery += "?,"
		db.Args = append(db.Args, val)
	}
	db.WhereSqlQuery = db.WhereSqlQuery[:len(db.WhereSqlQuery)-1]
	db.WhereSqlQuery += ") "
	return db
}
func (db *Query) WhereNotInRaw(field string, childQuery ChildQuery, val ...interface{}) *Query {
	if db.WhereSqlQuery != "" {
		db.WhereSqlQuery += "and "
	}
	db.WhereSqlQuery += "`" + field + "` not in ("
	check := &Query{}
	childQuery(check, val...)
	check.OperateFindToSql()
	db.WhereSqlQuery += check.SqlQuery
	for _, val := range check.Args {
		db.Args = append(db.Args, val)
	}
	db.WhereSqlQuery += ") "
	return db
}
func (db *Query) WhereBetween(field string, begin interface{}, over interface{}) *Query {
	if db.WhereSqlQuery != "" {
		db.WhereSqlQuery += "and ("
	}
	db.WhereSqlQuery += "`" + field + "` BETWEEN ? AND ?) "
	db.Args = append(db.Args, begin)
	db.Args = append(db.Args, over)
	return db
}
func (db *Query) WhereNotBetween(field string, begin interface{}, over interface{}) *Query {
	if db.WhereSqlQuery != "" {
		db.WhereSqlQuery += "and ("
	}
	db.WhereSqlQuery += "`" + field + "` NOT BETWEEN ? AND ?) "
	db.Args = append(db.Args, begin)
	db.Args = append(db.Args, over)
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
func (db *Query) Raw(field string, args ...interface{}) *Query {
	db.WhereSqlQuery += field
	for _, val := range args {
		db.Args = append(db.Args, val)
	}
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
		db.SqlQuery += StringBySliceString(",", db.RecordField) + " "
	} else {
		db.SqlQuery += "* "
	}
}
func (db *Query) OperateFindTable() {
	db.SqlQuery += "FROM " + db.RecordTable + " "
}
func (db *Query) OperateFindGroupBy() {
	db.SqlQuery += db.RecordGroup
}
func (db *Query) OperateFindOrderBy() {
	db.SqlQuery += db.RecordOrder
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

func (db *Query) InsertByMap(data *map[string]interface{}) (sql.Result, error) {
	db.OperateInsertTable()
	db.OperateInsertDataByMap(data)
	if db.Tx != nil {
		return db.Tx.Exec(db.SqlQuery, db.Args...)
	}
	return Open().Exec(db.SqlQuery, db.Args...)
}
func (db *Query) InsertByStruct(data interface{}) (sql.Result, error) {
	db.OperateInsertTable()
	db.OperateInsertDataByStruct(data)
	if db.Tx != nil {
		return db.Tx.Exec(db.SqlQuery, db.Args...)
	}
	return Open().Exec(db.SqlQuery, db.Args...)
}
func (db *Query) InsertAllByMap(datas *[]map[string]interface{}) (sql.Result, error) {
	db.OperateInsertTable()
	for key, val := range *datas {
		if key == 0 {
			db.OperateInsertDataByMap(&val)
		} else {
			db.OperateInsertDataByMapValue(&val)
		}
	}
	if db.Tx != nil {
		return db.Tx.Exec(db.SqlQuery, db.Args...)
	}
	return Open().Exec(db.SqlQuery, db.Args...)
}
func (db *Query) InsertAllByStruct(datas []interface{}) (sql.Result, error) {
	db.OperateInsertTable()
	for key, val := range datas {
		if key == 0 {
			db.OperateInsertDataByStruct(val)
		} else {
			db.OperateInsertDataByStructValue(val)
		}
	}
	if db.Tx != nil {
		return db.Tx.Exec(db.SqlQuery, db.Args...)
	}
	return Open().Exec(db.SqlQuery, db.Args...)

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
	fmt.Println(db.SqlQuery, db.Args)
}

// +----------------------------------------------------------------------
// | 更改方法
// +----------------------------------------------------------------------

func (db *Query) UpdateByMap(data *map[string]interface{}) (sql.Result, error) {
	db.OperateUpdateByMapData(data)
	if db.Tx != nil {
		return db.Tx.Exec(db.SqlQuery, db.Args...)
	}
	return Open().Exec(db.SqlQuery, db.Args...)
}
func (db *Query) UpdateByStruct(data interface{}) (sql.Result, error) {
	db.OperateUpdateByStructData(data)
	if db.Tx != nil {
		return db.Tx.Exec(db.SqlQuery, db.Args...)
	}
	return Open().Exec(db.SqlQuery, db.Args...)
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
func (db *Query) Delete() (sql.Result, error) {
	db.OperateDeleteData()
	if db.Tx != nil {
		return db.Tx.Exec(db.SqlQuery, db.Args...)
	}
	return Open().Exec(db.SqlQuery, db.Args...)
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

func (db *Query) Try(tx *sql.Tx) *Query {
	db.Tx = tx
	return db
}
