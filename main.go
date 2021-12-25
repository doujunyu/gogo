package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"gogo/ff"
	"gogo/sql"
	"io/ioutil"
	"os"
	"reflect"
	"time"
)

type Msg struct {
	Name    string
	Message string
	Number  int
}

type Hero struct {
	UserId int `json:"user_id"`
	CatId  int `json:"cat_id"`
}

func ReflectTypeValue(itf interface{}) {
	rtype := reflect.TypeOf(itf)
	fmt.Println("reflect type is ", rtype)
	rvalue := reflect.ValueOf(itf)
	fmt.Println("reflect value is ", rvalue)
	fmt.Println("reflect  value kind is", rvalue.Kind())
	fmt.Println("reflect type kind is", rtype.Kind())
	fmt.Println("reflect  value type is", rvalue.Type())
}

func ReflectStructElem(itf interface{}) {
	rvalue := reflect.ValueOf(itf)
	for i := 0; i < rvalue.NumField(); i++ {
		elevalue := rvalue.Field(i)
		fmt.Println("element ", i, " its type is ", elevalue.Type())
		fmt.Println("element ", i, " its kind is ", elevalue.Kind())
		fmt.Println("element ", i, " its value is ", elevalue)
	}
}

//type Demo struct {
//	LogChan *chan LogFunc
//}
type LogFunc func (string)
func Logs(set string){
	fmt.Println("我出来了:",set)
}
func ChanOut(a chan LogFunc){
	for true {
		data := <- a
		//data()
		fmt.Println(data)
	}
}
func ChanIn(a chan LogFunc,str string){
	a <- func(a string){
		Logs(a)
	}
}
//缓存测试

func main123() {
	//a := make(chan LogFunc,10)
	//a <- func("sdfs")

	//for {
	//	data := <- a
	//	//data()
	//	fmt.Println(data)
	//}
	//go ChanOut(a)
	//time.Sleep(time.Second *5)

	//redis
	//utility.RedisLine().Do("set","doudou", "sss")
	//data ,err :=utility.RedisLine().Do("get","doudou")
	//if err != nil {
	//	 fmt.Println(err)
	//	 return
	//}
	//fmt.Println(string(data.([]byte)),err)



	//defer conn.Close()
	//data,err := conn.Do("get","doudou")
	//fmt.Println(data,err)
	//fmt.Println(string(data.([]byte)),err)
	//sql 查询
	set := sql.Db("fs_users")
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
		fmt.Println(err)
		return
	}
	fmt.Println(data,err)

	//map类型添加
	//dataMap := make(map[string]interface{})
	//dataMap["user_id"] = 1
	//dataMap["cat_id"] = "123"
	//arr := make([]map[string]interface{},2)
	//arr[0] = dataMap
	//arr[1] = dataMap
	//data,err := sql.Db("sx_user_like").InsertAllByMap(&arr)
	//fmt.Println(data,err)
	//结构体添加
	//arr := make([]interface{},2)
	//arr[0] = &Hero{UserId: 1, CatId: 1}
	//arr[1] = &Hero{UserId: 2, CatId: 2}
	//data,err := sql.Db("sx_user_like").InsertAllByStruct(arr)
	//fmt.Println(data,err)
	//re ,err := sql.Exec("insert into sx_user_like(user_id, cat_id)values(?, ?)", "1", 1)
	//id, err :=re.LastInsertId()
	//fmt.Println(id,err)
	//map更改
	//dataMap := make(map[string]interface{})
	//dataMap["user_id"] = 1
	//dataMap["cat_id"] = "123"
	//data ,err := sql.Db("sx_user_like").WhereId("61454").UpdateByMap(&dataMap)
	//fmt.Println(data,err)
	//结构体更改
	//hh := Hero{UserId: 1, CatId: 1}
	//data,err :=sql.Db("sx_user_like").WhereId("61454").UpdateByStruct(&hh)
	//fmt.Println(data,err)
	//删除
	//set,err :=sql.Db("sx_user_like").WhereId("61454").Delete()
	//fmt.Println(set,err)

}
func midd1(jg *ff.Job) {
	fmt.Println("中间件执行了")
}

func main() {
	r := ff.ReadyGo()
	//r.Server.Handler = http.TimeoutHandler(http. NotFoundHandler(),time.Second*20,"sdfs")
	//r.Middleware = append(r.Middleware, midd1)
	//*r.LogChan <- 10
	r.GET("/demo", func(j *ff.Job) {

		j.W.Header().Set("Content-Type", "text/html")

		if j.Input["doudou"] == "1"{
			*r.ServerClose<-1
			j.JsonSuccess(nil,"正在关机",1)
			return
		}
		fmt.Println(j.Input["doudou"])
		j.Log.Error("ceshi","这是一个记录")
		//j.Log.Write("demo","doudou","ddddd")
		//j.Log.Error("doudou","sefsef")
		//time.Sleep(time.Second * 5)
		//_ = tx.Rollback()

		//j.JsonError(nil,"顺利通过",1)
		j.JsonSuccess(nil,"顺利通过",1)
	})


	r.GET("/file", func(j *ff.Job) {
		files := make(map[string]interface{})
		files["size"] = int64(1024 * 1024 * 5)
		files["name"] = "sss" + time.Now().Format("2006-01-02-15-04-05")
		files["suffix"] = "png,jpg"

		file, err := j.InputFile("file", "demo", files)
		if err != nil {
			j.JsonSuccess(file, "bbbb")
		}
		j.JsonSuccess(file, "aaaaa")
		fmt.Println(file, err)
	})

	r.GET("/public/", func(j *ff.Job) {
		file, _ := os.Open("." + j.R.URL.Path)
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
			}
		}(file)
		buff, _ := ioutil.ReadAll(file)
		_,_ = j.W.Write(buff)
	})


	go r.LogChanOut()
	go r.SetClose()
	_ = r.Server.ListenAndServe()
}

