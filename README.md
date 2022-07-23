# gogo

#### 介绍
针对http开发所使用的框架，模仿gin写的框架。
QQ群可讨论：6848027（GoLang/Go语言开发群）

#### 软件架构
内部通过简单的封装实现基本的http访问，请求http返回json数据，db数据库，redis,缓存，日志，软关闭服务


#### 安装教程

1.  go get github.com/doujunyu/gogo

#### 使用说明
一个简单的实例，更多操作请看同级main.go包

	r := gogo.ReadyGo()//生成配置
	r.GET("/demo",func(j *job.Job){//j.w 和 j.r 是http.ResponseWriter和*http.Request
        j.Log.Error("ceshi","这是一个记录")//记录日志
        j.JsonSuccess(j.Input,"这是信息",0)//返回输出数据,什么都不填会返回code = 0,
	})
	r.Run(":7070")


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
