package main

import (
	"database/sql"
	"fmt"
	"github.com/doujunyu/gogo/cache"
	"github.com/doujunyu/gogo/gogo"
	"github.com/doujunyu/gogo/gogo_log"
	"github.com/doujunyu/gogo/job"
	"github.com/doujunyu/gogo/sql_aid"

	//_ "github.com/go-sql-driver/mysql"//mysql数据库
	//_ "github.com/lib/pq" //pg数据库
	"io/ioutil"
	"os"
	"time"
)

var MySqlLine *sql.DB
var PGSqlLine *sql.DB
func main() {
	//链接数pgsql
	pgsql,err := sql_aid.Open("postgres",os.Getenv("PGSQL_URL"))
	fmt.Println(os.Getenv("PGSQL_URL"))
	if err != nil {
		fmt.Println("数据路链接错误",err)
		return
	}
	PGSqlLine = pgsql
	PGSqlLine.SetConnMaxLifetime(time.Minute * 3)
	PGSqlLine.SetMaxOpenConns(10)
	PGSqlLine.SetMaxIdleConns(10)
	//链接mysql
	mysql, err := sql.Open("mysql", os.Getenv("MYSQL_URL"))
	if err != nil {
		fmt.Println("数据路链接错误")
		return
	}
	MySqlLine = mysql
	MySqlLine.SetConnMaxLifetime(time.Minute * 3)
	MySqlLine.SetMaxOpenConns(10)
	MySqlLine.SetMaxIdleConns(10)




	r := gogo.ReadyGo()
	r.GET("/demosql", func(j *job.Job) {
		gogo_log.Write("前缀","gogo","正常信息")
		data := make(map[string]interface{})
		data["user_id"] = 1
		data["cat_id"] = 1

		set,slic := sql_aid.MyTable("sx_user_like").InsertByMap(&data) //生成sql语句
		tx,err := MySqlLine.Begin()                             //开启事务
		defer tx.Rollback()
		datas,err := tx.Exec(set,slic...)//进行添加
		lists := make([]int,1)//开始写让程序报错
		lists[0] = 1
		lists[2] = 1
		err = tx.Commit()

		fmt.Println(datas,err)
		fmt.Println(datas.LastInsertId())
		if err != nil{
			j.JsonError(nil,err.Error())
		}
		j.JsonSuccess()
	})

	//简单的例子
	r.GET("/demo", func(j *job.Job) {
		input := j.Input
		gogo_log.Error("demo","记录一条错误信息")
		//gogo_log.Write("前缀","gogo","正常信息")
		j.JsonSuccess(input,"这里是get提交")
	})
	r.POST("/demo", func(j *job.Job) {
		//input := j.Input
		gogo_log.Error("demo","记录一条错误信息")
		gogo_log.Write("前缀","gogo","正常信息")
		j.JsonSuccess(nil,"这里是post提交")
	})
	//之前中间件
	r.GET("/beforeGroup",group, func(j *job.Job) {
		input := j.Input
		gogo_log.Error("demo","记录一条错误信息")
		gogo_log.Write("前缀","gogo","正常信息")
		j.JsonSuccess(input)
	})
	//之后中间件
	r.GET("/laterGroup", func(j *job.Job) {
		input := j.Input
		gogo_log.Error("demo","记录一条错误信息")
		gogo_log.Write("前缀","gogo","正常信息")
		j.JsonSuccess(input)
	},group)
	//数据库查询
	//r.GET("/SqlFind", func(j *job.Job) {
	//	set := sql.Db("THIS_TABLE")
	//	set.Field("id", "nickname")
	//	set.WhereId("3")
	//	set.Where("openid", "like", "%4o1Bs%")
	//	set.Where("status","!=","1")
	//	set.WhereInRaw("id",func(child *sql.Query,val ...interface{}){
	//		child.Table("fs_user_address")
	//		child.Field("user_id","path")
	//		child.Where("status2","=",2)
	//	})
	//	set.WhereOrRaw(func(child *sql.Query,val ...interface{}){
	//		child.Where("status3","=",val[0])
	//		child.Where("status4","=",val[1])
	//		child.WhereBetween("status6",6,6.3)
	//		child.WhereOrRaw(func(child *sql.Query,val ...interface{}){
	//			child.Where("status5","=",5)
	//		})
	//	},3,4)
	//	set.OrderBy("id desc")
	//	set.PageSize(1,10)
	//	data, err := set.FindOnly()
	//	if err != nil {
	//		j.JsonError()
	//		return
	//	}
	//	j.JsonSuccess(data)
	//})
	////数据添加
	//r.GET("/SqlTryAdd",func(j *job.Job){
	//	//map类型添加
	//	dataMap := make(map[string]interface{})
	//	dataMap["user_id"] = 1
	//	dataMap["cat_id"] = "123"
	//	arr := make([]map[string]interface{},2)
	//	arr[0] = dataMap
	//	arr[1] = dataMap
	//	tx := sql.Try()
	//	data,err := sql.Db("sx_user_like").Try(tx).InsertAllByMap(&arr)
	//	if err != nil {
	//		j.JsonError(nil,err)
	//		tx.Rollback()
	//		return
	//	}
	//	tx.Commit()
	//	j.JsonSuccess(data)
	//})
	//缓存
	r.GET("/cache", func(j *job.Job) {
		if j.Input["data"] != "" {

			cache.Set(j.Input["name"],j.Input["data"],5)
		}
		data :=  cache.Get(j.Input["name"])
		if data == nil {
			j.JsonError(nil,"缓存过期")
			return
		}
		getTime := cache.GetTime(j.Input["name"])
		datas := []interface{}{data,getTime}
		j.JsonSuccess(datas)
	})
	//软关闭服务
	r.GET("/over", func(j *job.Job) {
		if j.Input["account"] == "gogo" && j.Input["password"] == "123456"{
			j.JsonError(nil,"账号或密码错误")
			return
		}
		//r.ServerClose<-1 //执行关机
		j.JsonSuccess(nil,"正在关机")
	})
	//文件上传
	r.GET("/file", func(j *job.Job) {
		files := make(map[string]interface{})
		files["size"] = int64(1024 * 1024 * 5)
		files["name"] = "sss" + time.Now().Format("2006-01-02-15-04-05")
		files["suffix"] = "png,jpg"
		file, err := j.InputFile("file", "demo", files)
		if err != nil {
			j.JsonError(nil, err)
		}
		j.JsonSuccess(file, "aaaaa")
		fmt.Println(file, err)
	})
	//访问public下面的文件
	r.GET("/public/", func(j *job.Job) {
		file, _ := os.Open("." + j.R.URL.Path)
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
			}
		}(file)
		buff, _ := ioutil.ReadAll(file)
		_,_ = j.W.Write(buff)
	})
	r.Run(":8000")

}

func group(j *job.Job){
	//做逻辑之后可以把需要的值通过这个方式传给后面要执行的逻辑方法中
	j.GroupData["token"] = 123456
}
