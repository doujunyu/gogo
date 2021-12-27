package main

import (
	"fmt"
	"github.com/doujunyu/gogo/cache"
	"github.com/doujunyu/gogo/gogo"
	"github.com/doujunyu/gogo/job"
	"github.com/doujunyu/gogo/sql"
	"io/ioutil"
	"os"
	"time"
)



func main() {
	r := gogo.ReadyGo()
	//简单的例子
	r.GET("/demo", func(j *job.Job) {
		input := j.Input
		j.Log.Error("demo","记录一条错误信息")
		j.Log.Write("前缀","gogo","正常信息")
		j.JsonSuccess(input)
	})
	//之前中间件
	r.GET("/beforeGroup",group, func(j *job.Job) {
		input := j.Input
		j.Log.Error("demo","记录一条错误信息")
		j.Log.Write("前缀","gogo","正常信息")
		j.JsonSuccess(input)
	})
	//之后中间件
	r.GET("/laterGroup", func(j *job.Job) {
		input := j.Input
		j.Log.Error("demo","记录一条错误信息")
		j.Log.Write("前缀","gogo","正常信息")
		j.JsonSuccess(input)
	},group)
	//数据库查询
	r.GET("/SqlFind", func(j *job.Job) {
		set := sql.Db("THIS_TABLE")
		set.Field("id", "nickname")
		set.WhereId("3")
		set.Where("openid", "like", "%4o1Bs%")
		set.Where("status","!=","1")
		set.WhereInRaw("id",func(child *sql.Query,val ...interface{}){
			child.Table("fs_user_address")
			child.Field("user_id","path")
			child.Where("status2","=",2)
		})
		set.WhereOrRaw(func(child *sql.Query,val ...interface{}){
			child.Where("status3","=",val[0])
			child.Where("status4","=",val[1])
			child.WhereBetween("status6",6,6.3)
			child.WhereOrRaw(func(child *sql.Query,val ...interface{}){
				child.Where("status5","=",5)
			})
		},3,4)
		set.OrderBy("id desc")
		set.PageSize(1,10)
		data, err := set.FindOnly()
		if err != nil {
			j.JsonError()
			return
		}
		j.JsonSuccess(data)
	})
	//缓存
	r.GET("/cache", func(j *job.Job) {
		cache.Set(j.Input["name"],j.Input["data"],30)
		cat := cache.Get(j.Input["name"])
		j.JsonSuccess(cat)
	})
	//软关闭服务
	r.GET("/over", func(j *job.Job) {
		if j.Input["account"] == "gogo" && j.Input["password"] == "123456"{
			j.JsonError(nil,"账号或密码错误")
			return
		}
		*r.ServerClose<-1 //执行关机
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

	r.Run(":7070")
}

func group(j *job.Job){
	//做逻辑之后可以把需要的值通过这个方式传给后面要执行的逻辑方法中
	j.GroupData["token"] = 123456
}
