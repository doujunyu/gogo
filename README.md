# gogo

#### 介绍
针对http开发所使用的框架，模仿gin写的框架。

#### 软件架构
内部通过简单的封装实现基本的http访问，请求http返回json数据，db数据库，redis,缓存，日志，软关闭服务


#### 安装教程

1.  go get github.com/doujunyu/gogo

#### 使用说明
简单示例

	r := gogo.ReadyGo()//生成配置
	r.GET("/demo",func(j *gogo.Job){//j.w 和 j.r 是http.ResponseWriter和*http.Request
        j.Log.Error("ceshi","这是一个记录")//记录日志
        j.JsonSuccess(j.Input,"这是信息",0)//返回输出数据,什么都不填会返回code = 0,
	})
	go r.LogChanOut()
	go r.SetClose()
	_ = r.Server.ListenAndServe()
数据库操作

    //数据查询
    r.GET("/demo",func(j *gogo.Job){
    set := gogo.Db("THIS_TABLE")
    set.Field("id", "nickname")
    set.WhereId("3")
    set.Where("openid", "like", "%4o1Bs%")
    set.Where("status","!=","1")
    set.WhereInRaw("id",func(child *gogo.Query,val ...interface{}){
        child.Table("fs_user_address")
        child.Field("user_id","path")
        child.Where("status2","=",2)
    })
    set.WhereOrRaw(func(child *gogo.Query,val ...interface{}){
        child.Where("status3","=",val[0])
        child.Where("status4","=",val[1])
        child.WhereBetween("status6",6,6.3)
        child.WhereOrRaw(func(child *gogo.Query,val ...interface{}){
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
---
    //map添加
    dataMap := make(map[string]interface{})
    dataMap["user_id"] = 1
    dataMap["cat_id"] = "123"
    data,err := gogo.Db("sx_user_like").InsertByMap(&dataMap)
    //多个添加用切片包起来使用InsertAllByMap
---
    //结构体添加
    type Hero struct {
    UserId int `json:"user_id"`
    CatId  int `json:"cat_id"`
    }

    gogo.Db("sx_user_like").InsertByStruct(&Hero{UserId: 1, CatId: 1})
    //多个添加用切片包起来使用InsertAllByStruct
---
    //map更改
    dataMap := make(map[string]interface{})
    dataMap["user_id"] = 1
    dataMap["cat_id"] = "123"
    data ,err := gogo.Db("sx_user_like").WhereId("61454").UpdateByMap(&dataMap)
---
    //结构体更改
    hh := Hero{UserId: 1, CatId: 1}
    data,err :=gogo.Db("sx_user_like").WhereId("61454").UpdateByStruct(&hh)
---
    //删除
    set,err :=gogo.Db("sx_user_like").WhereId("61454").Delete()
    fmt.Println(set,err)
---

缓存操作


#### 参与贡献

1.  Fork 本仓库
2.  新建 Feat_xxx 分支
3.  提交代码
4.  新建 Pull Request


#### 特技

1.  使用 Readme\_XXX.md 来支持不同的语言，例如 Readme\_en.md, Readme\_zh.md
2.  Gitee 官方博客 [blog.gitee.com](https://blog.gitee.com)
3.  你可以 [https://gitee.com/explore](https://gitee.com/explore) 这个地址来了解 Gitee 上的优秀开源项目
4.  [GVP](https://gitee.com/gvp) 全称是 Gitee 最有价值开源项目，是综合评定出的优秀开源项目
5.  Gitee 官方提供的使用手册 [https://gitee.com/help](https://gitee.com/help)
6.  Gitee 封面人物是一档用来展示 Gitee 会员风采的栏目 [https://gitee.com/gitee-stars/](https://gitee.com/gitee-stars/)
