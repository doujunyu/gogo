package cache

import (
	"sync"
	"time"
)

var GlobalCacheData *cache
var GlobalCacheOnce sync.Once

//写入
type cache struct {
	CacheChan chan int         `Testing:"生成管道,添加和编辑用管道操作"`
	Data      map[string]*cacheValue `Testing:"操作的数据"`
}

// CacheValue 缓存数据格式
type cacheValue struct {
	Data interface{}
	Time int64
}

func init()  {
	GlobalCacheOnce.Do(func() {
		GlobalCacheData =  &cache{
			CacheChan: make(chan int,1),
			Data:      make(map[string]*cacheValue),
		}
	})
}

// Set 新增缓存,时间=0表示长期存储,无论有没有数据都直接赋值
func Set(name string, data interface{}, second int64) bool {
	GlobalCacheData.CacheChan <- 1
	GlobalCacheData.Data[name] = &cacheValue{
		Time: time.Now().Unix() + second,
		Data: data,
	}
	<-GlobalCacheData.CacheChan
	return true
}

// SetCover Set 新增缓存,时间=0表示长期存储,如果数据存在则不会进行覆盖,并返回false
func SetCover(name string, data interface{}, second int64) bool {
	GlobalCacheData.CacheChan <- 1
	if !Exists(name) {
		<-GlobalCacheData.CacheChan
		return false
	}
	GlobalCacheData.Data[name] = &cacheValue{
		Time: time.Now().Unix() + second,
		Data: data,
	}
	<-GlobalCacheData.CacheChan
	return true
}

// SetDataCover 重置缓存数据
func SetDataCover(name string, data interface{}) bool {
	GlobalCacheData.CacheChan <- 1
	if !Exists(name) {
		<-GlobalCacheData.CacheChan
		return false
	}
	GlobalCacheData.Data[name].Data = data
	<-GlobalCacheData.CacheChan
	return true
}

// SetTimeCover 重置缓存时间
func SetTimeCover(name string, second int64) bool {
	GlobalCacheData.CacheChan <- 1
	if !Exists(name) {
		<-GlobalCacheData.CacheChan
		return false
	}
	GlobalCacheData.Data[name].Time = time.Now().Unix() + second
	<-GlobalCacheData.CacheChan
	return true
}

// SetTimeIncrease 增加缓存时间
func SetTimeIncrease(name string, second int64) bool {
	GlobalCacheData.CacheChan <- 1
	if !Exists(name) {
		<-GlobalCacheData.CacheChan
		return false
	}
	GlobalCacheData.Data[name].Time +=  second
	<-GlobalCacheData.CacheChan
	return true
}

// SetTimeDecrease 减少缓存时间
func SetTimeDecrease(name string, second int64) bool {
	GlobalCacheData.CacheChan <- 1
	if !Exists(name) {
		<-GlobalCacheData.CacheChan
		return false
	}
	GlobalCacheData.Data[name].Time -=  second
	<-GlobalCacheData.CacheChan
	return true
}

// Del 删除
func Del(name string) bool {
	GlobalCacheData.CacheChan <- 1
	if !Exists(name) {
		<-GlobalCacheData.CacheChan
		return true
	}
	delete(GlobalCacheData.Data, name)
	<-GlobalCacheData.CacheChan
	return true
}

// Get 获取缓存数据
func Get(name string) interface{} {
	if !Exists(name) {
		return nil
	}
	return GlobalCacheData.Data[name].Data
}

// GetTime 获取缓存时间
func GetTime(name string) int64 {
	if !Exists(name) {
		return 0
	}
	return GlobalCacheData.Data[name].Time - time.Now().Unix()
}

// Exists 判断数据是否存在
func Exists(name string) bool {
	if GlobalCacheData.Data[name] == nil || time.Now().Unix() > GlobalCacheData.Data[name].Time {
		return false
	}
	return true
}

// ChanLongTime 清理过期的缓存
func ChanLongTime() {
	for {
		time.Sleep(time.Second)
		for name, value := range GlobalCacheData.Data {
			time.Sleep(time.Second)
			if value.Time < time.Now().Unix() {
				Del(name)
			}
		}
	}
}




