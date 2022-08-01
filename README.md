- [gogo](#gogo)
    - [介绍](#介绍)
        - [软件架构](#软件架构)
    - [安装：](#安装)
    - [使用说明：](#使用说明)
# gogo
#### 介绍
针对http开发所使用的框架，模仿gin写的框架。
QQ群可讨论：6848027（GoLang/Go语言开发群）

#### 软件架构
内部通过简单的封装实现基本的:
    http访问简化，请求http返回json数据，db数据库工厂模式，redis封装,缓存，日志，杀死进程软关闭


#### 安装
1. go get github.com/doujunyu/gogo
2. go get gitee.com/doujunyu/gogo

#### 使用说明
一个简单的实例，更多操作请看同级main.go包

	r := gogo.ReadyGo()//生成配置
	r.GET("/demo",func(j *job.Job){//j.w 和 j.r 是http.ResponseWriter和*http.Request
        j.Log.Error("ceshi","这是一个记录")//记录日志
        j.JsonSuccess(j.Input,"这是信息",0)//返回输出数据,什么都不填会返回code = 0,
	})
	r.Run(":7070")
